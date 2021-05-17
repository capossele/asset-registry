package main

import (
	"flag"
	"sync"
	"time"

	"github.com/capossele/asset-registry/pkg/registryservice"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	// defaultMongoDBOpTimeout defines the default MongoDB operation timeout.
	defaultMongoDBOpTimeout = 5 * time.Second
)

var (
	clientDB   *mongo.Client
	dbOnce     sync.Once
	server     *echo.Echo
	serverOnce sync.Once
	log        *zap.SugaredLogger
)

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log = logger.Sugar()

	service := registryservice.NewService(mongoDB())
	httpHandler := registryservice.NewHTTPHandler(service, log)

	Server()

	// configure the server
	server.HideBanner = true
	server.HidePort = true

	server.GET("/", IndexRequest)
	server.POST("/registries/:network/assets", httpHandler.SaveAsset)
	server.GET("/registries/:network/assets", httpHandler.LoadAssets)
	server.GET("/registries/:network/assets/:ID", httpHandler.LoadAsset)

	log.Infof("Starting server ...")

	log.Fatal(server.Start(*httpBindAddr))
}
