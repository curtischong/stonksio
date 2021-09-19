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
		Host:        "https://nitter.net/search/rss",
		Query:       "f=tweets&q=ethereum&e-media=on&e-images=on&e-videos=on&e-native_video=on&e-pro_video=on&e-replies=on&e-nativeretweets=on",
		AverageTime: 10,
	}, ch)

	feed.Start()

	for post := range ch {
		fmt.Println(post.Body)
	}
}
