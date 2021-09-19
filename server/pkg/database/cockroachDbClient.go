package database

import (
	"context"
	"fmt"
	"log"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"time"

	"github.com/google/uuid"

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
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}

	connConfig.Database = "stonksio"
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

func (client *CockroachDbClient) InsertPost(
	post common.Post,
) error {
	return crdbpgx.ExecuteTx(context.Background(), client.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		log.Printf("Creating post=%s\n", post)
		_, err := tx.Exec(context.Background(),
			`INSERT INTO post (id, username, "userpicurl", "body", "timestamp") VALUES ($1, $2, $3, $4, $5)`,
			uuid.New().String(), post.Username, post.UserPicUrl, post.Body, post.Timestamp)
		return err
	})
}

func (client *CockroachDbClient) deleteAllPosts() error {
	return crdbpgx.ExecuteTx(context.Background(), client.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		log.Printf("Deleting all posts")
		_, err := tx.Exec(context.Background(), "DELETE FROM post")
		return err
	})
}

func (client *CockroachDbClient) GetPosts(n int) ([]common.Post, error) {
	rows, err := client.conn.Query(context.Background(), `SELECT 'id', 'username', 'userpicurl', 'body', 'timestamp' FROM post LIMIT $1;`, n)
	if err != nil {
		return nil, fmt.Errorf("cannot query rows. err=%s", err)
	}
	posts := make([]common.Post, 0, n)
	defer rows.Close()
	for rows.Next() {
		post := common.Post{}
		var timestamp string
		if err := rows.Scan(&post.Id, &post.Username, &post.UserPicUrl, &post.Body, &timestamp); err != nil {
			return nil, fmt.Errorf("cannot scan rows. err=%s", err)
		}
		post.Timestamp, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return nil, fmt.Errorf("cannot parse post.Timestamp=%s, err=%s", timestamp, err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (client *CockroachDbClient) GetPrices(
	asset string,
) ([]common.Price, error) {
	if asset != "ETH" {
		return nil, fmt.Errorf("invalid asset=%s", asset)
	}
	rows, err := client.conn.Query(context.Background(), "SELECT tradePrice, timestamp FROM price WHERE asset=$1", asset)
	if err != nil {
		return nil, err
	}
	prices := make([]common.Price, 0)
	defer rows.Close()
	for rows.Next() {
		price := common.Price{}
		var timestamp string
		if err := rows.Scan(&price.TradePrice, &timestamp); err != nil {
			return nil, err
		}
		price.Timestamp, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return nil, fmt.Errorf("cannot parse price.Timestamp=%s, err=%s", timestamp, err)
		}
		prices = append(prices, price)
	}
	return prices, nil
}

func (client *CockroachDbClient) GetLatestPrice(
	asset string,
) (float32, error) {
	if asset != "ETH" {
		return 0, fmt.Errorf("invalid asset=%s", asset)
	}

	rows, err := client.conn.Query(context.Background(),
		"SELECT tradePrice FROM price WHERE timestamp = MAX(timestamp)")
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		break
	}
	return 0, nil
}

func (client *CockroachDbClient) InsertPrice(
	asset string, tradePrice float32,
) error {
	if asset != "ETH" {
		return fmt.Errorf("invalid asset=%s", asset)
	}
	return crdbpgx.ExecuteTx(context.Background(), client.conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		log.Printf("Creating tradePrice=%f for asset=%s\n", tradePrice, asset)
		_, err := tx.Exec(context.Background(),
			"INSERT INTO tradePrice (asset, tradePrice, timestamp) VALUES ($1, $2, $3)",
			asset, tradePrice, time.Now())
		return err
	})
}
