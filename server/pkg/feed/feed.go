package feed

import (
	"context"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"time"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

const (
	username              = "ape"
	scraperBackoffSeconds = 5
	query                 = "ethereum -filter:media -filter:retweets"
)

type Feed struct {
	logger      *log.Logger
	averageTime int
	scraper     *twitterscraper.Scraper
	out         chan<- *common.Post
	buf         chan *common.Post
}

func NewFeed(config config.FeedConfig, postChannel chan<- *common.Post) *Feed {
	scraper := twitterscraper.New()
	scraper.SetSearchMode(twitterscraper.SearchTop)
	scraper.WithDelay(scraperBackoffSeconds)

	return &Feed{
		logger:      log.New(),
		averageTime: config.AverageTime,
		scraper:     twitterscraper.New(),
		out:         postChannel,
		buf:         make(chan *common.Post, 50),
	}
}

func (f *Feed) Start() {
	// spawn fetch job
	go f.fetch()

	// spawn producer job
	go f.producer()
}

func (f *Feed) producer() {
	defer func() {
		if err := recover(); err != nil {
			f.logger.Warnf("Feed producer job died, restarting. err=%s", err)
			go f.producer()
		}
	}()

	for {
		delay := time.Duration(rand.Intn(2*f.averageTime)) * time.Second
		time.Sleep(delay)

		post := <-f.buf
		post.Timestamp = time.Now()
		f.out <- post
	}
}

func (f *Feed) fetch() {
	defer func() {
		if err := recover(); err != nil {
			f.logger.Warnf("Feed fetch job died, restarting. err=%s", err)
			go f.fetch()
		}
	}()

	for {
		for tweet := range f.scraper.SearchTweets(context.Background(), query, 50) {
			f.buf <- &common.Post{
				Id:       tweet.ID,
				Username: username, // not actual tweet username
				Body:     tweet.Text,
			}
		}
	}
}
