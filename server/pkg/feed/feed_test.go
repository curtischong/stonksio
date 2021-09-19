package feed_test

import (
	"fmt"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	feed2 "stonksio/pkg/feed"
	"testing"
)

func TestFeed(t *testing.T) {
	ch := make(chan *common.Post)
	feed := feed2.NewFeed(config.FeedConfig{
		AverageTime: 10,
	}, ch)

	feed.Start()

	for post := range ch {
		fmt.Println(post.Body)
	}
}
