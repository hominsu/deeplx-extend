package main

import (
	"flag"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/oschwald/geoip2-golang"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/oio-network/deeplx-extend/app/deeplx/internal/conf"
	"github.com/oio-network/deeplx-extend/app/deeplx/internal/task"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "deeplx-extend.deeplx"

	// Version is the version of the compiled software.
	Version string

	// flagconf is the config flag.
	flagconf string

	// flagGeoLite is the IP database
	flagGeoLite string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&flagGeoLite, "geolite2", "../../data/GeoLite2-Country.mmdb", "MaxMind's GeoLite2 Country database, eg: -geolite2 GeoLite2-Country.mmdb")

	json.MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: false,
	}
}

func newApp(logger log.Logger, ms task.MachineryServer, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			ms,
			hs,
		),
	)
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer func(c config.Config) {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}(c)
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	db, err := geoip2.Open(flagGeoLite)
	if err != nil {
		panic(err)
	}
	defer func(db *geoip2.Reader) {
		if err := db.Close(); err != nil {
			log.Error(err)
		}
	}(db)

	logger := log.With(log.NewStdLogger(os.Stdout),
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"ts", log.Timestamp(time.RFC3339),
		"caller", log.Caller(6),
	)

	app, cleanup, err := initApp(bc.Server, bc.Data, bc.Secret, db, logger, &fasthttp.Client{})
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
