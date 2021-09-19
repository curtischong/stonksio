package ohlc

import (
	"stonksio/pkg/common"
	"stonksio/pkg/database"
	"stonksio/pkg/sentiment"

	log "github.com/sirupsen/logrus"
)

const defaultLastPrice = 3000.0

type OhlcGenerator struct {
	cockroachDbClient *database.CockroachDbClient
	gcpClient         sentiment.GcpClient
	lastPrice         float32
}

func NewOhlcGenerator(
	cockroachDbClient *database.CockroachDbClient,
	gcpClient sentiment.GcpClient,
	asset string,
) *OhlcGenerator {
	lastPrice, err := cockroachDbClient.GetLatestOhlc(asset)
	if err != nil {
		log.Errorf("couldn't find last price. err=%s", err)
		lastPrice = defaultLastPrice
	}

	return &OhlcGenerator{
		cockroachDbClient: cockroachDbClient,
		gcpClient:         gcpClient,
		lastPrice:         lastPrice,
	}
}

func (generator *OhlcGenerator) writeNewPrice(
	post common.Post,
) error {
	_, err := generator.gcpClient.CalculateSentiment(post.Body)
	if err != nil {
		return err
	}

	return nil
}
