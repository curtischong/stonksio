package main

import (
	"net/http"
	"stonksio/pkg/common"
	"stonksio/pkg/conductor"
	"stonksio/pkg/database"
	"stonksio/pkg/feed"
	"stonksio/pkg/generator/price"
	"stonksio/pkg/post"
	"stonksio/pkg/request"
	"stonksio/pkg/sentiment"
	"stonksio/pkg/websocket"

	log "github.com/sirupsen/logrus"

	"stonksio/pkg/config"
)

const configPath = "./config.yaml"

func main() {
	config, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("couldn't load config path=%s, err=%s", configPath, err)
	}

	cockroachDbClient := database.NewCockroachDbClient(config)
	gcpClient := sentiment.NewGcpClient()
	pusherClient := websocket.NewPusherClient()

	incomingPrices := make(chan *common.Price)
	priceGenerator := price.NewPriceGenerator(cockroachDbClient, gcpClient, incomingPrices, "ETH")

	postHandler := post.NewPostHandler(cockroachDbClient, priceGenerator, pusherClient)

	incomingPosts := make(chan *common.Post)
	feedSrv := feed.NewFeed(config.Feed, incomingPosts)

	conductorSrv := conductor.NewConductor(cockroachDbClient, postHandler, pusherClient, incomingPosts, incomingPrices)

	requestHandler := request.NewRequestHandler(config, cockroachDbClient, postHandler)

	http.HandleFunc("/api/post", requestHandler.HandlePostPost)
	http.HandleFunc("/api/posts", requestHandler.HandleGetPosts)
	http.HandleFunc("/api/prices/eth", requestHandler.HandleGetPrices)
	http.HandleFunc("/api/ohlc/eth", requestHandler.HandleGetOHLCs)
	http.HandleFunc("/api/wallet", requestHandler.HandleGetWallet)

	// start
	priceGenerator.Start()
	feedSrv.Start()
	conductorSrv.Start()

	log.Info("Starting server on port 8090")
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatalf("Cannot start server err=%s", err)
	}
}
