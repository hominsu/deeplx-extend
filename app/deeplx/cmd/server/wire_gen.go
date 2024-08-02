// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/biz"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/conf"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/data"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/server"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/service"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/task"
	"github.com/oio-network/deeplx-extend/app/deeplx/pkgs/client_pool"
	"github.com/oio-network/deeplx-extend/deeplx"
	"github.com/oschwald/geoip2-golang"
	"github.com/valyala/fasthttp"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confData *conf.Data, confSecret *conf.Secret, db *geoip2.Reader, logger log.Logger, clients ...*fasthttp.Client) (*kratos.App, func(), error) {
	client := data.NewEntClient(confData, logger)
	cmdable := data.NewRedisCmd(confData, logger)
	cacheAsideClient, cleanup := data.NewCacheAsideClient(confData, logger)
	cache := data.NewCache(cacheAsideClient)
	dataData, cleanup2, err := data.NewData(client, cmdable, cache, confData, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	accessLogRepo := data.NewAccessLogRepo(dataData, logger)
	accessLogUseCase := biz.NewAccessLogUsecase(accessLogRepo, logger)
	logSink, cleanup3 := task.NewLogSink(confData, accessLogUseCase, logger)
	logTask := task.NewLogTask(db, logSink, logger)
	machineryServer := task.NewMachineryServer(confData, logTask)
	translateService := deeplx.NewTranslateService(logger)
	clientPool, cleanup4 := client_pool.NewClientPool(clients...)
	deepLXService := service.NewDeepLXService(machineryServer, translateService, confSecret, clientPool, logger)
	httpServer := server.NewHTTPServer(confServer, deepLXService, logger)
	app := newApp(logger, machineryServer, httpServer)
	return app, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
