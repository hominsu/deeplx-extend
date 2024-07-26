package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/cache/v9"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"

	"github.com/oio-network/deeplx-extend/app/deeplx/internal/conf"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/data/ent"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/data/ent/migrate"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewEntClient,
	NewRedisCmd,
	NewRedisCache,
	NewUserRepo,
	NewAccessLogRepo,
)

type Data struct {
	db    *ent.Client
	rdCmd redis.Cmdable
	cache *cache.Cache

	conf *conf.Data
}

func NewData(
	entClient *ent.Client,
	rdCmd redis.Cmdable,
	conf *conf.Data,
	logger log.Logger,
) (*Data, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "data"))

	data := &Data{
		db:    entClient,
		rdCmd: rdCmd,
		conf:  conf,
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
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           int(conf.Redis.Db),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()

	err := client.Ping(timeout).Err()
	if err != nil {
		helper.Fatalf("redis connect error: %v", err)
	}
	return client
}

func NewRedisCache(rdCmd redis.Cmdable, conf *conf.Data) *cache.Cache {
	opts := &cache.Options{
		Redis: rdCmd,
	}
	if conf.Cache.LfuEnable {
		opts.LocalCache = cache.NewTinyLFU(int(conf.Cache.LfuSize), conf.Cache.Ttl.AsDuration())
	}

	return cache.New(opts)
}
