package websocket

import (
	pusher "github.com/pusher/pusher-http-go"
	"stonksio/pkg/common"
)

// the file can't be called pusherClient to for syntax highlighting to work for some reason
type PusherClient struct {
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
		socketClient: socketClient,
	}
}

func (client *PusherClient) PushPost(
	post *common.Post,
) {
	client.socketClient.Trigger("post", "new-post", post)
}

func (client *PusherClient) PushPrice(
	price *common.Price,
) {
	client.socketClient.Trigger("prices", "new-price", price)
}
