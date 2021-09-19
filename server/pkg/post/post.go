package post

import (
	"stonksio/pkg/common"
	"stonksio/pkg/database"
	"stonksio/pkg/generator/price"
	"stonksio/pkg/websocket"
)

type PostHandler struct {
	cockroachDbClient *database.CockroachDbClient
	priceGenerator    *price.PriceGenerator
	pusherClient      *websocket.PusherClient
}

func NewPostHandler(
	cockroachDbClient *database.CockroachDbClient,
	priceGenerator *price.PriceGenerator,
	pusherClient *websocket.PusherClient,
) *PostHandler {
	return &PostHandler{
		cockroachDbClient: cockroachDbClient,
		priceGenerator:    priceGenerator,
		pusherClient:      pusherClient,
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

	return nil
}
