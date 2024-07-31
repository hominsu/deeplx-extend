package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	v1 "github.com/oio-network/deeplx-extend/api/deeplx/v1"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/conf"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/task"
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

	lt *task.LogTask
	ts *deeplx.TranslateService

	cs   *conf.Secret
	pool *client_pool.ClientPool

	log *log.Helper
}

func NewDeepLXService(
	lt *task.LogTask,
	ts *deeplx.TranslateService,
	cs *conf.Secret,
	pool *client_pool.ClientPool,
	logger log.Logger,
) *DeepLXService {
	return &DeepLXService{
		lt:   lt,
		ts:   ts,
		cs:   cs,
		pool: pool,
		log:  log.NewHelper(log.With(logger, "module", "service/deeplx")),
	}
}
