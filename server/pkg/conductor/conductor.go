package conductor

import (
	"stonksio/pkg/common"
	"stonksio/pkg/database"
	"stonksio/pkg/post"
	"stonksio/pkg/websocket"

	log "github.com/sirupsen/logrus"
)

type Conductor struct {
	cockroachDbClient *database.CockroachDbClient
	postHandler       *post.PostHandler
	pusherClient      *websocket.PusherClient
	incomingPosts     <-chan *common.Post
	incomingPrices    <-chan *common.Price
	logger            *log.Logger
}

func NewConductor(
	cockroachDbClient *database.CockroachDbClient,
	postHandler *post.PostHandler,
	pusherClient *websocket.PusherClient,
	incomingPosts <-chan *common.Post,
	incomingPrices <-chan *common.Price,
) *Conductor {
	conductor := &Conductor{
		logger:            log.New(),
		cockroachDbClient: cockroachDbClient,
		postHandler:       postHandler,
		pusherClient:      pusherClient,
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
			if err := c.cockroachDbClient.InsertPrice(*price); err != nil {
				log.Errorf("cannot insert price err=%s", err)
			} else {
				c.pusherClient.PushPrice(price)
			}
		}
	}
}
