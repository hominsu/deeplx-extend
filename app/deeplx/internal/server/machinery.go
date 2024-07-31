package server

import (
	"context"

	"github.com/oio-network/deeplx-extend/app/deeplx/internal/conf"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/task"
	"github.com/oio-network/deeplx-extend/pkgs/machinery"
)

const LogTaskCreateAccessLog = "LogTask.CreateAccessLog"

type LogTask interface {
	CreateAccessLog(remoteAddr string) error
}

func RegisterLogTask(srv MachineryServer, t LogTask) error {
	if err := srv.HandleFunc(LogTaskCreateAccessLog, t.CreateAccessLog); err != nil {
		return err
	}

	return nil
}

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
	t *task.LogTask,
) MachineryServer {
	opts := []machinery.ServerOption{
		machinery.WithBrokerAddress(c.Redis.Addr, int(c.Redis.Db), machinery.BrokerTypeRedis),
		machinery.WithResultBackendAddress(c.Redis.Addr, int(c.Redis.Db), machinery.BackendTypeRedis),
	}

	srv := machinery.NewServer(opts...)

	if err := RegisterLogTask(srv, t); err != nil {
		panic(err)
	}

	return srv
}
