package database

import (
	"context"
	"log"
	"stonksio/pkg/common"
	"stonksio/pkg/config"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"

	"github.com/jackc/pgx/v4"
)

type CockroachDbClient struct {
	config config.DatabaseConfig
	conn   *pgx.Conn
}

func NewCockroachDbClient(
	config *config.Config,
) *CockroachDbClient {
	// Connect to the stonksio database
	connConfig, err := pgx.ParseConfig(config.DatabaseConfig.ConnectionString)
	connConfig.Database = "stonksio"
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	return &CockroachDbClient{
		config: config.DatabaseConfig,
		conn:   conn,
	}
}

func (client *CockroachDbClient) closeConn() {
	client.conn.Close(context.Background())
}

func (client *CockroachDbClient) insertPost(
	post common.Post,
) error {
	return crdbpgx.ExecuteTx(context.Background(), client.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		log.Printf("Creating post=%s\n", post)
		_, err := tx.Exec(context.Background(),
			"INSERT INTO post (id, message) VALUES ($1, $2, $3)", post.Username, post.UserPicUrl, post.Body)
		return err
	})
}

func (client *CockroachDbClient) deleteallposts() error {
	return crdbpgx.ExecuteTx(context.Background(), client.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		log.Printf("Deleting all posts")
		_, err := tx.Exec(context.Background(), "DELETE FROM post")
		return err
	})
}
