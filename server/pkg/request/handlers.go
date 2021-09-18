package request

import (
	"encoding/json"
	"net/http"
	"stonksio/pkg/config"
	"stonksio/pkg/database"
)

type RequestHandler struct {
	config            *config.Config
	cockroachDbClient *database.CockroachDbClient
}

func NewRequestHandler(
	config *config.Config,
) *RequestHandler {
	return &RequestHandler{
		config:            config,
		cockroachDbClient: database.NewCockroachDbClient(config),
	}
}

func (handler *RequestHandler) HandleGetOhlc(
	w http.ResponseWriter, r *http.Request,
) {
	handler.sendStatusOK(w)
	ohlc, err := handler.cockroachDbClient.GetOhlc("ETH")
	if err != nil {
		handler.sendInternalServerError(w, err)
	}
	handler.sendStatusOK(w)
	fileUrlsBytes, _ := json.Marshal(ohlc)
	w.Write(fileUrlsBytes)
}
