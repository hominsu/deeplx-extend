package task

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/oschwald/geoip2-golang"

	"github.com/oio-network/deeplx-extend/app/deeplx/internal/conf"
	"github.com/oio-network/deeplx-extend/pkgs/machinery"
)

var ProviderSet = wire.NewSet(
	NewLogTask,
	NewLogSink,
	NewMachineryServer,
)

type MachineryServer interface {
	HandleFunc(name string, handler any) error

	NewTask(typeName string, opts ...machinery.TaskOption) error
	NewPeriodicTask(cronSpec, typeName string, opts ...machinery.TaskOption) error
	NewGroup(groupTasks ...machinery.TasksOption) error
	NewPeriodicGroup(cronSpec string, groupTasks ...machinery.TasksOption) error
	NewChord(chordTasks ...machinery.TasksOption) error
	NewPeriodicChord(cronSpec string, chordTasks ...machinery.TasksOption) error
	NewChain(chainTasks ...machinery.TasksOption) error
	NewPeriodicChain(cronSpec string, chainTasks ...machinery.TasksOption) error

	Start(context.Context) error
	Stop(context.Context) error
}

func NewMachineryServer(
	c *conf.Data,
	t LogTask,
) MachineryServer {
	opts := []machinery.ServerOption{
		machinery.WithBrokerAddress(c.Redis.Addr, int(c.Redis.Db), machinery.BrokerTypeRedis),
		machinery.WithResultBackendAddress(c.Redis.Addr, int(c.Redis.Db), machinery.BackendTypeRedis),
	}

	srv := machinery.NewServer(opts...)

	if err := t.RegisterLogTask(srv); err != nil {
		panic(err)
	}

	return srv
}

type logTask struct {
	mmdb *geoip2.Reader
	sink *LogSink

	log *log.Helper
}

type LogTask interface {
	RegisterLogTask(srv MachineryServer) error

	CreateAccessLog(b []byte) error
}

func NewLogTask(
	mmdb *geoip2.Reader,
	sink *LogSink,
	logger log.Logger,
) LogTask {
	return &logTask{
		mmdb: mmdb,
		sink: sink,
		log:  log.NewHelper(log.With(logger, "module", "task/log")),
	}
}
