package task

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/oschwald/geoip2-golang"

	"github.com/oio-network/deeplx-extend/app/deeplx/internal/biz"
)

var ProviderSet = wire.NewSet(
	NewLogTask,
)

type LogTask struct {
	au   *biz.AccessLogUseCase
	mmdb *geoip2.Reader
	log  *log.Helper
}

func NewLogTask(
	au *biz.AccessLogUseCase,
	mmdb *geoip2.Reader,
	logger log.Logger,
) *LogTask {
	return &LogTask{
		au:   au,
		mmdb: mmdb,
		log:  log.NewHelper(log.With(logger, "module", "task/log")),
	}
}
