package database

import (
	"fmt"
	"stonksio/pkg/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewDb(t *testing.T) {
	config, err := config.NewConfig("../../config.yaml")
	assert.NoError(t, err)

	dbClient := NewCockroachDbClient(config)
	t.Run("viewPosts", func(t *testing.T) {
		posts, err := dbClient.GetPosts(1)
		fmt.Println(posts)
		assert.NoError(t, err)
	})
}

func TestInsertDb(t *testing.T) {
	config, err := config.NewConfig("../../config.yaml")
	assert.NoError(t, err)

	dbClient := NewCockroachDbClient(config)
	t.Run("createPost", func(t *testing.T) {
		posts, err := dbClient.GetPosts(1)
		fmt.Println(posts)
		assert.NoError(t, err)
	})
}
