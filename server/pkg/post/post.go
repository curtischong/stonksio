package post

import (
	"fmt"
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

func (h *PostHandler) HandlePost(post *common.Post) {
	fmt.Println(post)
}
