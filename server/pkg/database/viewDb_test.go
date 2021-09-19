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
		posts, err := dbClient.GetPosts(1)
		fmt.Println(posts)
		assert.NoError(t, err)
	})

	t.Run("viewPosts", func(t *testing.T) {
		posts, err := dbClient.GetPosts(1)
		assert.NoError(t, err)
		print(posts)
	})
}

func TestInsertDb(t *testing.T) {
	config, err := config.NewConfig("../../config.yaml")
	assert.NoError(t, err)

	dbClient := NewCockroachDbClient(config)
	t.Run("viewPosts", func(t *testing.T) {
		err := dbClient.InsertPost(common.Post{
			uuid.New().String(),
			"asdfasd",
			"asdfasd",
			"asdfasd",
			time.Now(),
		})
		assert.NoError(t, err)
	})
}

func TestDeleteDb(t *testing.T) {
	_, err := config.NewConfig("../../config.yaml")
	assert.NoError(t, err)

	//dbClient := NewCockroachDbClient(config)
	/*t.Run("deletePost", func(t *testing.T) {
		posts, err := dbClient.delete(1)
		fmt.Println(posts)
		assert.NoError(t, err)
	})*/
}
