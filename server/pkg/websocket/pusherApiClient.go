package websocket

import (
	pusher "github.com/pusher/pusher-http-go"
	"stonksio/pkg/common"

	log "github.com/sirupsen/logrus"
)

// the file can't be called pusherClient to for syntax highlighting to work for some reason
type PusherClient struct {
	logger       *log.Logger
	socketClient pusher.Client
}

func NewPusherClient() *PusherClient {
	socketClient := pusher.Client{
		AppID:   "1269168",
		Key:     "f710317ee72763936d91",
		Secret:  "dbc764b648a553776edc",
		Cluster: "us2",
		Secure:  true,
	}
	return &PusherClient{
		logger:       log.New(),
		socketClient: socketClient,
	}
}

func (client *PusherClient) PushPost(
	post *common.Post,
) {
	client.logger.Infof("pushing post %s", post.Id)
	if err := client.socketClient.Trigger("post", "new-post", post); err != nil {
		client.logger.Errorf("could not push post, err=%s", err)
	}
}

func (client *PusherClient) PushPrice(
	price *common.Price,
) {
	client.logger.Infof("pushing price %f", price.TradePrice)
	if err := client.socketClient.Trigger("prices", "new-price", price); err != nil {
		client.logger.Errorf("could not push price, err=%s", err)
	}
}
