package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
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
	pool *client_pool.ClientPool
	log  *log.Helper
}

func NewDeepLXService(ts *deeplx.TranslateService, pool *client_pool.ClientPool, logger log.Logger) *DeepLXService {
	return &DeepLXService{
		ts:   ts,
		pool: pool,
		log:  log.NewHelper(log.With(logger, "module", "service/deeplx")),
	}
}
