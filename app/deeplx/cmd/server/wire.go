//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/oschwald/geoip2-golang"
	"github.com/valyala/fasthttp"

	"github.com/oio-network/deeplx-extend/app/deeplx/internal/conf"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/server"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/service"
)

// initApp init kratos application.
func initApp(
	confServer *conf.Server,
	confSecret *conf.Secret,
	db *geoip2.Reader,
	logger log.Logger,
	clients ...*fasthttp.Client,
) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			service.ProviderSet,
			server.ProviderSet,
			newApp,
		),
	)
}
