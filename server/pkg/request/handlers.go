package request

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"stonksio/pkg/database"
	"stonksio/pkg/post"
	"stonksio/pkg/websocket"

	log "github.com/sirupsen/logrus"
)

type RequestHandler struct {
	logger            *log.Logger
	config            *config.Config
	cockroachDbClient *database.CockroachDbClient
	PusherClient      *websocket.PusherClient
	postHandler       *post.PostHandler
}

func NewRequestHandler(
	config *config.Config,
	cockroachDbClient *database.CockroachDbClient,
	postHandler *post.PostHandler,
) *RequestHandler {
	return &RequestHandler{
		logger:            log.New(),
		config:            config,
		cockroachDbClient: cockroachDbClient,
		PusherClient:      websocket.NewPusherClient(),
		postHandler:       postHandler,
	}
}

func (handler *RequestHandler) HandleGetPrices(
	w http.ResponseWriter, r *http.Request,
) {
	prices, err := handler.cockroachDbClient.GetPrices("ETH")
	if err != nil {
		handler.sendInternalServerError(w, err)
	}
	handler.sendStatusOK(w)
	fileUrlsBytes, _ := json.Marshal(prices)
	w.Write(fileUrlsBytes)
}

func (handler *RequestHandler) HandlePostPost(
	w http.ResponseWriter, r *http.Request,
) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.logger.Errorf("error reading the body err=%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error", "message": "internal server error"}`))
		return
	}

	var newPost common.Post
	if err := json.Unmarshal(body, &newPost); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "error", "message": "could not parse json"}`))
		return
	}

	// Create ID for post
	newPost.Id = uuid.New().String()

	if err := handler.postHandler.HandlePost(&newPost); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error", "message": "internal server error"}`))
		return
	}

	w.Write([]byte(`{"status": "success"}`))
}
