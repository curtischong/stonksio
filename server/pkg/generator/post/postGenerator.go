package post

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"time"

	log "github.com/sirupsen/logrus"
)

type PostGenerator struct {
	tweets map[string][]string // maps tweeters -> tweets
}

func newPostGenerator(
	config *config.Config,
) *PostGenerator {
	tweets := make(map[string][]string, 0)
	for _, tweetFile := range config.DatasetConfig.TweetFiles {
		jsonFile, err := os.Open(config.DatasetConfig.DatasetRoot + "/" + tweetFile)
		if err != nil {
			log.Fatalf("cannot open tweetFile=%s, err=%s", tweetFile, err)
		}
		defer jsonFile.Close()
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatalf("cannot read tweetFile=%s, err=%s", tweetFile, err)
		}
		currentTweets := make([]string, 0)
		json.Unmarshal(byteValue, &currentTweets)

		tweeter := tweetFile[:len(tweetFile)-5]
		tweets[tweeter] = currentTweets
	}
	return &PostGenerator{
		tweets: tweets,
	}
}

func (generator *PostGenerator) nextPost(
	tweeter string,
) common.Post {
	rand.Seed(time.Now().UnixNano())
	tweetPool := generator.tweets[tweeter]
	min := 0
	max := len(tweetPool)
	randIdx := rand.Intn(max-min+1) + min
	return common.Post{
		Username: tweeter,
		Body:     tweetPool[randIdx],
	}
}
