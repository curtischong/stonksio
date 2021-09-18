package main

import (
	"net/http"
	"stonksio/pkg/request"

	log "github.com/sirupsen/logrus"

	"stonksio/pkg/config"
)

const configPath = "./config.yaml"

func main() {
	config, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("couldn't load config path=%s, err=%s", configPath, err)
	}

	requestHandler := request.NewRequestHandler(config)
	http.HandleFunc("/get/ohlc/eth", requestHandler.HandleGetOhlc)
	log.Info("Starting server on port 8090")
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatalf("Cannot start server err=%s", err)
	}
}
