package request

import (
	"encoding/json"
	io "io/ioutil"
	"net/http"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"stonksio/pkg/database"
	"stonksio/pkg/post"
	"stonksio/pkg/websocket"
	"strconv"

	"github.com/google/uuid"

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
	var num = 50
	if n := r.URL.Query().Get("n"); n != "" {
		var err error
		num, err = strconv.Atoi(n)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "could not parse 'n'",
			})
			return
		}
	}
	prices, err := handler.cockroachDbClient.GetPrices("ETH", num)
	if err != nil {
		handler.sendInternalServerError(w, err)
	}
	handler.sendStatusOK(w)
	json.NewEncoder(w).Encode(prices)
}

func (handler *RequestHandler) HandlePostPost(
	w http.ResponseWriter, r *http.Request,
) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.logger.Errorf("error reading the body err=%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "internal server error",
		})
		return
	}

	var newPost common.Post
	if err := json.Unmarshal(body, &newPost); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "could not parse json",
		})
		return
	}

	// Create ID for post
	newPost.Id = uuid.New().String()

	if err := handler.postHandler.HandlePost(&newPost); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}
