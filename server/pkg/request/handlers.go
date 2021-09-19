package request

import (
	"encoding/json"
	"fmt"
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

func (handler *RequestHandler) HandleGetWallet(
	w http.ResponseWriter, r *http.Request,
) {
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "please send a username",
		})
		return
	}

	wallet, err := handler.cockroachDbClient.GetWallet("ETH", username)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}
	if wallet == nil {
		// create balance for user
		err := handler.cockroachDbClient.InsertWallet(common.Wallet{
			Username: username,
			Asset:    "USD",
			Balance:  10000,
		})
		if err != nil {
			handler.sendInternalServerError(w, err)
			return
		}
		err = handler.cockroachDbClient.InsertWallet(common.Wallet{
			Username: username,
			Asset:    "ETH",
			Balance:  50,
		})
		if err != nil {
			handler.sendInternalServerError(w, err)
			return
		}
	}

	handler.sendStatusOK(w)

	ethWallet, err := handler.cockroachDbClient.GetWallet("ETH", username)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}
	usdWallet, err := handler.cockroachDbClient.GetWallet("USD", username)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}
	walletData := make(map[string]string)
	walletData["ETH"] = fmt.Sprintf("%f", ethWallet.Balance)
	walletData["USD"] = fmt.Sprintf("%f", usdWallet.Balance)
	fileUrlsBytes, _ := json.Marshal(walletData)
	w.Write(fileUrlsBytes)
}

func (handler *RequestHandler) HandleBuy(
	w http.ResponseWriter, r *http.Request,
) {
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "please specify a username",
		})
		return
	}
	sizeStr := r.URL.Query().Get("orderSize")
	orderSize, err := strconv.Atoi(sizeStr)
	if err != nil || orderSize <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": fmt.Sprintf("invalid orderSize=%s specified. It needs to be a positive int", sizeStr),
		})
		return
	}

	// check if they have enough money
	usdWallet, err := handler.cockroachDbClient.GetWallet("USD", username)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}
	ethPrice, err := handler.cockroachDbClient.GetLatestPrice("ETH")
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	orderPrice := float32(orderSize) * ethPrice
	if orderPrice < usdWallet.Balance {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "insufficient USD",
		})
		return
	}
	usdWallet.Balance -= orderPrice
	handler.cockroachDbClient.UpdateWallet(*usdWallet)

	ethWallet, err := handler.cockroachDbClient.GetWallet("ETH", username)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}
	ethWallet.Balance += float32(orderSize)
	handler.cockroachDbClient.UpdateWallet(*ethWallet)
}

func (handler *RequestHandler) HandleSell(
	w http.ResponseWriter, r *http.Request,
) {
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "please specify a username",
		})
		return
	}
	sizeStr := r.URL.Query().Get("orderSize")
	orderSize, err := strconv.Atoi(sizeStr)
	if err != nil || orderSize <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": fmt.Sprintf("invalid orderSize=%s specified. It needs to be a positive int", sizeStr),
		})
		return
	}

	// check if they have enough ETH
	ethWallet, err := handler.cockroachDbClient.GetWallet("ETH", username)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	if orderSize > int(ethWallet.Balance) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "insufficient ETH",
		})
		return
	}

	ethPrice, err := handler.cockroachDbClient.GetLatestPrice("ETH")
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	usdWallet, err := handler.cockroachDbClient.GetWallet("USD", username)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	usdWallet.Balance += ethPrice * float32(orderSize)
	handler.cockroachDbClient.UpdateWallet(*usdWallet)

	ethWallet.Balance -= float32(orderSize)
	handler.cockroachDbClient.UpdateWallet(*ethWallet)
}

func (handler *RequestHandler) HandleGetOHLCs(
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

	ohlcs, err := handler.cockroachDbClient.GetOHLCs(count)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}
	handler.sendStatusOK(w)
	json.NewEncoder(w).Encode(ohlcs)
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
