package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/oschwald/geoip2-golang"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/conf"
	"github.com/oio-network/deeplx-extend/app/deeplx/pkgs/client_pool"
	"github.com/oio-network/deeplx-extend/deeplx"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewDeepLXService,
	client_pool.NewClientPool,
	deeplx.NewTranslateService,
)

type DeepLXService struct {
	v1.UnimplementedDeepLXServiceServer

	ts   *deeplx.TranslateService
	cs   *conf.Secret
	pool *client_pool.ClientPool
	mmdb *geoip2.Reader
	log  *log.Helper
}

func NewDeepLXService(
	ts *deeplx.TranslateService,
	cs *conf.Secret,
	pool *client_pool.ClientPool,
	mmdb *geoip2.Reader,
	logger log.Logger,
) *DeepLXService {
	return &DeepLXService{
		ts:   ts,
		cs:   cs,
		pool: pool,
		mmdb: mmdb,
		log:  log.NewHelper(log.With(logger, "module", "service/deeplx")),
	}
}
