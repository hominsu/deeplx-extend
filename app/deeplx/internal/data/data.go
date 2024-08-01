package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisaside"

	"github.com/oio-network/deeplx-extend/app/deeplx/internal/conf"
	"github.com/oio-network/deeplx-extend/app/deeplx/pkgs/msg"
	"github.com/oio-network/deeplx-extend/schema/ent"
	"github.com/oio-network/deeplx-extend/schema/ent/migrate"

	// driver
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewEntClient,
	NewRedisCmd,
	NewCacheAsideClient,
	NewCache,
	NewUserRepo,
	NewAccessLogRepo,
)

type Data struct {
	db    *ent.Client
	rdCmd redis.Cmdable
	cache *Cache

	conf *conf.Data
	log  *log.Helper
}

func NewData(
	entClient *ent.Client,
	rdCmd redis.Cmdable,
	cache *Cache,
	conf *conf.Data,
	logger log.Logger,
) (*Data, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "data"))

	data := &Data{
		db:    entClient,
		rdCmd: rdCmd,
		cache: cache,
		conf:  conf,
		log:   helper,
	}
	return data, func() {
		if err := data.db.Close(); err != nil {
			helper.Error(err)
		}
	}, nil
}

func NewEntClient(conf *conf.Data, logger log.Logger) *ent.Client {
	helper := log.NewHelper(log.With(logger, "module", "data/ent"))

	client, err := ent.Open(
		conf.Database.Driver,
		conf.Database.Source,
	)
	if err != nil {
		helper.Fatalf("failed opening connection to db: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(
		context.Background(),
		migrate.WithForeignKeys(true),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		helper.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}

func NewRedisCmd(conf *conf.Data, logger log.Logger) redis.Cmdable {
	helper := log.NewHelper(log.With(logger, "module", "data/redis"))

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       int(conf.Redis.Db),
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()

	err := client.Ping(timeout).Err()
	if err != nil {
		helper.Fatalf("redis connect error: %v", err)
	}
	return client
}

func NewCacheAsideClient(conf *conf.Data, logger log.Logger) (rueidisaside.CacheAsideClient, func()) {
	helper := log.NewHelper(log.With(logger, "module", "data/redis-cache-aside"))

	cc, err := rueidisaside.NewClient(rueidisaside.ClientOption{
		ClientOption: rueidis.ClientOption{
			InitAddress:  []string{conf.Redis.Addr},
			Password:     conf.Redis.Password,
			SelectDB:     int(conf.Redis.Db),
			DisableCache: conf.Cache.DisableClientCache,
		},
	})
	if err != nil {
		helper.Fatalf("redis connect error: %v", err)
	}

	return cc, func() {
		cc.Close()
	}
}

type (
	MarshalFunc   func(any) ([]byte, error)
	UnmarshalFunc func([]byte, any) error
)

type Cache struct {
	client    rueidisaside.CacheAsideClient
	Marshal   MarshalFunc
	Unmarshal UnmarshalFunc
}

func NewCache(client rueidisaside.CacheAsideClient) *Cache {
	return &Cache{
		client:    client,
		Marshal:   msg.Marshal,
		Unmarshal: msg.Unmarshal,
	}
}

func (c *Cache) Get(ctx context.Context, ttl time.Duration, key string, value any, fn func(ctx context.Context) (value any, err error)) error {
	res, err := c.client.Get(ctx, ttl, key, func(ctx context.Context, key string) (val string, err error) {
		ret, err := fn(ctx)
		if err != nil {
			return "", err
		}

		bytes, err := c.Marshal(ret)
		if err != nil {
			return "", err
		}

		return string(bytes), nil
	})
	if err != nil {
		return err
	}

	return c.Unmarshal([]byte(res), value)
}
