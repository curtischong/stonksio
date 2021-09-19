package websocket

import (
	"stonksio/pkg/common"
	"testing"
	"time"
)

func TestPushPost(t *testing.T) {
	pusherClient := NewPusherClient()

	pusherClient.PushPost(&common.Post{
		"asdfasd",
		"splacorn",
		"https://avatars.githubusercontent.com/u/10677873?v=4",
		"I love Golang!",
		time.Now(),
	})
}
