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
	"time"

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

	incomingPrices := make(chan *common.Price)
	priceGenerator := price.NewPriceGenerator(cockroachDbClient, gcpClient, incomingPrices, "ETH")

	postHandler := post.NewPostHandler(cockroachDbClient, priceGenerator)

	incomingPosts := make(chan *common.Post)
	feedSrv := feed.NewFeed(config.Feed, incomingPosts)

	conductorSrv := conductor.NewConductor(config, cockroachDbClient, postHandler, incomingPosts, incomingPrices)

	requestHandler := request.NewRequestHandler(config, cockroachDbClient)

	sendTestPush(requestHandler)
	http.HandleFunc("api/post", requestHandler.HandlePostPost)
	http.HandleFunc("/api/prices/eth", requestHandler.HandleGetPrices)
	log.Info("Starting server on port 8090")

	// start
	priceGenerator.Start()
	feedSrv.Start()
	conductorSrv.Start()

	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatalf("Cannot start server err=%s", err)
	}
}

func sendTestPush(
	requestHandler *request.RequestHandler,
) {
	requestHandler.PusherClient.PushPost(common.Post{
		"asdfasd",
		"splacorn",
		"https://avatars.githubusercontent.com/u/10677873?v=4",
		"I love Golang!",
		time.Now(),
	})
}
