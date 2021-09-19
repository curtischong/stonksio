package websocket

import (
	"fmt"
	"stonksio/pkg/common"
	"time"

	pusher "github.com/pusher/pusher-http-go"
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
	post common.Post,
) {
	data := map[string]string{
		"id":         post.Id,
		"username":   post.Username,
		"userPicUrl": post.UserPicUrl,
		"body":       post.Body,
		"timestamp":  post.Timestamp.Format(time.RFC3339),
	}
	client.socketClient.Trigger("post", "new-post", data)
}

// TODO: call this function when updating prices
func (client *PusherClient) PushPrices(
	prices map[string]common.Price, // maps asset -> price
) {
	data := make(map[string]map[string]string)
	for asset, price := range prices {
		data[asset] = map[string]string{
			"tradePrice": fmt.Sprintf("%f", price.TradePrice),
			"timestamp":  price.Timestamp.Format(time.RFC3339),
		}
	}
	client.socketClient.Trigger("prices", "new-prices", data)
}
