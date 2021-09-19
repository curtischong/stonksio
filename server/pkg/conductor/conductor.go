package conductor

import (
	"fmt"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"stonksio/pkg/database"
	"stonksio/pkg/feed"
	"stonksio/pkg/generator/price"
	"stonksio/pkg/sentiment"

	log "github.com/sirupsen/logrus"
)

type Conductor struct {
	feed              *feed.Feed
	priceGenerator    *price.PriceGenerator
	incomingPosts     chan *common.Post
	incomingPrices    chan *common.Price
	cockroachDbClient *database.CockroachDbClient
	logger            *log.Logger
}

func NewConductor(
	config *config.Config,
	cockroachDbClient *database.CockroachDbClient,
	gcpClient sentiment.GcpClient,
) *Conductor {
	incomingPosts := make(chan *common.Post)
	incomingPrices := make(chan *common.Price)
	conductor := &Conductor{
		logger:            log.New(),
		feed:              feed.NewFeed(config.Feed, incomingPosts),
		priceGenerator:    price.NewPriceGenerator(cockroachDbClient, gcpClient, incomingPrices, "ETH"),
		cockroachDbClient: cockroachDbClient,
		incomingPosts:     incomingPosts,
	}
	conductor.feed.Start()
	conductor.priceGenerator.Start()
	go conductor.Start()
	return conductor
}

func (c *Conductor) Start() {
	defer func() {
		if err := recover(); err != nil {
			c.logger.Warnf("Conductor job died, restarting. err=%s", err)
			go c.Start()
		}
	}()

	for {
		select {
		case post := <-c.incomingPosts:
			fmt.Println(post)
			// TODO: write post to db
			tradePrice, err := c.priceGenerator.GetNewPriceFromPostSentiment(post.Body)
			if err != nil {
				log.Errorf("cannot get price from post sentiment err=%s", err)
				continue
			}
			if err := c.cockroachDbClient.InsertPrice("ETH", tradePrice); err != nil {
				log.Errorf("cannot insert price err=%s", err)
				continue
			}

		case price := <-c.incomingPrices:
			fmt.Println(price)
			if err := c.cockroachDbClient.InsertPrice("ETH", price.TradePrice); err != nil {
				log.Errorf("cannot insert price err=%s", err)
				continue
			}
		}

	}
}
