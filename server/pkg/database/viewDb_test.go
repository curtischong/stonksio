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
		posts, err := dbClient.getPosts(1)
		fmt.Println(posts)
		assert.NoError(t, err)
	})

	t.Run("viewPosts", func(t *testing.T) {
		posts, err := dbClient.getPosts(1)
		assert.NoError(t, err)
		print(posts)
	})
}

func TestInsertDb(t *testing.T) {
	config, err := config.NewConfig("../../config.yaml")
	assert.NoError(t, err)

	dbClient := NewCockroachDbClient(config)
	t.Run("createPost", func(t *testing.T) {
		posts, err := dbClient.getPosts(1)
		fmt.Println(posts)
		assert.NoError(t, err)
	})
}
