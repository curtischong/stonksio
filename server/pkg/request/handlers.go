package request

import (
	"encoding/json"
	"net/http"
	"stonksio/pkg/config"
	"stonksio/pkg/database"
	"stonksio/pkg/websocket"
)

type RequestHandler struct {
	config            *config.Config
	cockroachDbClient *database.CockroachDbClient
	PusherClient      *websocket.PusherClient
}

func NewRequestHandler(
	config *config.Config,
) *RequestHandler {
	return &RequestHandler{
		config:            config,
		cockroachDbClient: database.NewCockroachDbClient(config),
		PusherClient:      websocket.NewPusherClient(),
	}
}

func (handler *RequestHandler) HandleGetPrices(
	w http.ResponseWriter, r *http.Request,
) {
	handler.sendStatusOK(w)
	prices, err := handler.cockroachDbClient.GetPrices("ETH")
	if err != nil {
		handler.sendInternalServerError(w, err)
	}
	handler.sendStatusOK(w)
	fileUrlsBytes, _ := json.Marshal(prices)
	w.Write(fileUrlsBytes)
}
