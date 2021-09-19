package request

import (
	"encoding/json"
	io "io/ioutil"
	"net/http"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"stonksio/pkg/database"
	"stonksio/pkg/post"
	"strconv"
	"time"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

const (
	defaultPricesWindow = 5 * time.Minute
	defaultPostsCount   = 20
)

type RequestHandler struct {
	logger            *log.Logger
	config            *config.Config
	cockroachDbClient *database.CockroachDbClient
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
		postHandler:       postHandler,
	}
}

func (handler *RequestHandler) HandleGetPrices(
	w http.ResponseWriter, r *http.Request,
) {
	window, err := getWindow(r, defaultPricesWindow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "could not parse window",
		})
		return
	}

	prices, err := handler.cockroachDbClient.GetPrices("ETH", window)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
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

func (handler *RequestHandler) HandleGetPosts(
	w http.ResponseWriter, r *http.Request,
) {
	count, err := getCount(r, defaultPostsCount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "could not parse count",
		})
		return
	}

	posts, err := handler.cockroachDbClient.GetPosts(count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func getCount(r *http.Request, defaultCount int) (int, error) {
	if count := r.URL.Query().Get("count"); count != "" {
		return strconv.Atoi(count)
	} else {
		return defaultCount, nil
	}
}

func getWindow(r *http.Request, defaultWindow time.Duration) (time.Duration, error) {
	if window := r.URL.Query().Get("window"); window != "" {
		return time.ParseDuration(window)
	} else {
		return defaultWindow, nil
	}
}
