package post

import (
	"stonksio/pkg/common"
	"stonksio/pkg/database"
	"stonksio/pkg/generator/price"
	"stonksio/pkg/ohlc"
	"stonksio/pkg/websocket"
)

type PostHandler struct {
	cockroachDbClient *database.CockroachDbClient
	priceGenerator    *price.PriceGenerator
	pusherClient      *websocket.PusherClient
	ohlcManager       *ohlc.OHLCManager
}

func NewPostHandler(
	cockroachDbClient *database.CockroachDbClient,
	priceGenerator *price.PriceGenerator,
	pusherClient *websocket.PusherClient,
	ohlcManager *ohlc.OHLCManager,
) *PostHandler {
	return &PostHandler{
		cockroachDbClient: cockroachDbClient,
		priceGenerator:    priceGenerator,
		pusherClient:      pusherClient,
		ohlcManager:       ohlcManager,
	}
}

func (h *PostHandler) HandlePost(post *common.Post) error {
	if err := h.cockroachDbClient.InsertPost(*post); err != nil {
		return err
	}

	newPrice, err := h.priceGenerator.GetNewPriceFromPostSentiment(post.Body)
	if err != nil {
		return err
	}
	if err := h.cockroachDbClient.InsertPrice(*newPrice); err != nil {
		return err
	}

	h.pusherClient.PushPost(post)
	h.pusherClient.PushPrice(newPrice)
	h.ohlcManager.HandlePrice(newPrice)

	return nil
}
