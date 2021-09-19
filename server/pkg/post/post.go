package post

import (
	"stonksio/pkg/common"
	"stonksio/pkg/database"
	"stonksio/pkg/generator/price"
)

type PostHandler struct {
	cockroachDbClient *database.CockroachDbClient
	priceGenerator    *price.PriceGenerator
}

func NewPostHandler(
	cockroachDbClient *database.CockroachDbClient,
	priceGenerator *price.PriceGenerator,
) *PostHandler {
	return &PostHandler{
		cockroachDbClient: cockroachDbClient,
		priceGenerator:    priceGenerator,
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
	if err := h.cockroachDbClient.InsertPrice("ETH", newPrice); err != nil {
		return err
	}

	return nil
}
