package price_test

import (
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"stonksio/pkg/database"
	"stonksio/pkg/generator/price"
	"stonksio/pkg/sentiment"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

// The purpose of this file is to make sure the generated prices make sense

func TestViewDb(t *testing.T) {
	config, err := config.NewConfig("../../../config.yaml")
	assert.NoError(t, err)

	dbClient := database.NewCockroachDbClient(config)
	gcpClient := sentiment.NewGcpClient()

	incomingPrices := make(chan *common.Price)
	generator := price.NewPriceGenerator(dbClient, gcpClient, incomingPrices, "ETH")
	t.Run("startGenerator", func(t *testing.T) {
		generator.Start()
		consumer(incomingPrices)
	})
}

func consumer(pchan chan *common.Price) {
	for p := range pchan {
		log.Infof("price: %f", p.TradePrice)
	}
}
