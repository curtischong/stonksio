package conductor

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"stonksio/pkg/database"
	"stonksio/pkg/post"
)

type Conductor struct {
	postHandler       *post.PostHandler
	incomingPosts     <-chan *common.Post
	incomingPrices    <-chan *common.Price
	cockroachDbClient *database.CockroachDbClient
	logger            *log.Logger
}

func NewConductor(
	config *config.Config,
	cockroachDbClient *database.CockroachDbClient,
	postHandler *post.PostHandler,
	incomingPosts <-chan *common.Post,
	incomingPrices <-chan *common.Price,
) *Conductor {
	conductor := &Conductor{
		logger:            log.New(),
		cockroachDbClient: cockroachDbClient,
		postHandler:       postHandler,
		incomingPosts:     incomingPosts,
		incomingPrices:    incomingPrices,
	}
	return conductor
}

func (c *Conductor) Start() {
	go c.consumer()
}

func (c *Conductor) consumer() {
	defer func() {
		if err := recover(); err != nil {
			c.logger.Warnf("Conductor job died, restarting. err=%s", err)
			go c.consumer()
		}
	}()

	for {
		select {
		case post := <-c.incomingPosts:
			c.postHandler.HandlePost(post)

		case price := <-c.incomingPrices:
			fmt.Println(price)
			if err := c.cockroachDbClient.InsertPrice("ETH", price.TradePrice); err != nil {
				log.Errorf("cannot insert price err=%s", err)
			}
		}
	}
}
