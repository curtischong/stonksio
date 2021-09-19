package main

import (
	"net/http"
	"stonksio/pkg/common"
	"stonksio/pkg/feed"
	"stonksio/pkg/request"
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

	_ = feed.NewFeed(config.Feed, nil)

	requestHandler := request.NewRequestHandler(config)

	sendTestPush(requestHandler)
	http.HandleFunc("/get/prices/eth", requestHandler.HandleGetPrices)
	http.HandleFunc("/post/post", requestHandler.HandlePostPost)
	log.Info("Starting server on port 8090")
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
