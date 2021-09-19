package database

import (
	"fmt"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestViewDb(t *testing.T) {
	config, err := config.NewConfig("../../config.yaml")
	assert.NoError(t, err)

	dbClient := NewCockroachDbClient(config)
	t.Run("viewPosts", func(t *testing.T) {
		posts, err := dbClient.GetPosts(100)
		fmt.Println(posts)
		assert.NoError(t, err)
	})

	t.Run("viewPrices", func(t *testing.T) {
		prices, err := dbClient.GetPrices("ETH", 5*time.Minute)
		fmt.Println(prices)
		assert.NoError(t, err)
	})

	t.Run("viewOHLCs", func(t *testing.T) {
		ohlcs, err := dbClient.GetOHLCs(5)
		fmt.Println(ohlcs)
		assert.NoError(t, err)
	})

	t.Run("viewBalance", func(t *testing.T) {
		prices, err := dbClient.GetBalance("ETH", "splacorn")
		fmt.Println(prices)
		assert.NoError(t, err)
	})
}

func TestInsertDb(t *testing.T) {
	config, err := config.NewConfig("../../config.yaml")
	assert.NoError(t, err)

	dbClient := NewCockroachDbClient(config)
	t.Run("insertPosts", func(t *testing.T) {
		err := dbClient.InsertPost(common.Post{
			uuid.New().String(),
			"asdfasd",
			"asdfasd",
			"asdfasd",
			time.Now(),
		})
		assert.NoError(t, err)
	})

	t.Run("insertPrice", func(t *testing.T) {
		err := dbClient.InsertPrice(common.Price{
			Asset:      "ETH",
			TradePrice: 2311.3,
			Timestamp:  time.Now(),
		})
		assert.NoError(t, err)
	})
}

func TestDeleteDb(t *testing.T) {
	config, err := config.NewConfig("../../config.yaml")
	assert.NoError(t, err)

	dbClient := NewCockroachDbClient(config)
	t.Run("deleteAllPosts", func(t *testing.T) {
		err := dbClient.deleteAllPosts()
		assert.NoError(t, err)
	})
	t.Run("deleteAllPrices", func(t *testing.T) {
		err := dbClient.deleteAllPrices()
		assert.NoError(t, err)
	})
}
