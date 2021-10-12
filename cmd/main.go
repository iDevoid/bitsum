package main

import (
	"flag"
	"os"

	"github.com/iDevoid/bitsum/internal/glue/routing"
	"github.com/iDevoid/bitsum/internal/handler/rest"
	"github.com/iDevoid/bitsum/internal/module/coins"
	"github.com/iDevoid/bitsum/internal/storage/memcache"
	"github.com/iDevoid/bitsum/platform/routers"
	"github.com/sirupsen/logrus"
)

var testInit bool

func init() {
	flag.BoolVar(&testInit, "test", false, "initialize test mode without serving")
	flag.Parse()
}

func main() {
	cache := memcache.Initialize()
	usecase := coins.Initialize(cache)
	handlers := rest.CoinsInit(usecase)
	routes := routing.CoinsRouting(handlers)

	server := routers.Initialize(":9000", routes, "coin")
	if testInit {
		logrus.Info("Initialize test mode Finished!")
		os.Exit(0)
	}

	server.Serve()
}
