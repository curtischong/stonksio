package database

import (
	"context"
	"fmt"
	"log"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"time"

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
			"INSERT INTO post (id, message) VALUES ($1, $2, $3, $4)", post.Username, post.UserPicUrl, post.Body, post.Timestamp)
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

func (client *CockroachDbClient) GetOhlc(
	asset string,
) ([]common.Ohlc, error) {
	if asset != "ETH" {
		return nil, fmt.Errorf("invalid asset=%s", asset)
	}
	rows, err := client.conn.Query(context.Background(), "SELECT open, high, low, close, startTime, endTime FROM ohlc")
	if err != nil {
		log.Fatal(err)
	}
	prices := make([]common.Ohlc, 0)
	defer rows.Close()
	for rows.Next() {
		ohlc := common.Ohlc{}
		var startTime string
		var endTime string
		if err := rows.Scan(&ohlc.Open, &ohlc.High, &ohlc.Low, &ohlc.Close, &startTime, &endTime); err != nil {
			return nil, err
		}
		ohlc.StartTime, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			return nil, fmt.Errorf("cannot parse ohlc.StartTime=%s, err=%s", startTime, err)
		}
		ohlc.EndTime, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			return nil, fmt.Errorf("cannot parse ohlc.EndTime=%s, err=%s", endTime, err)
		}
		prices = append(prices, ohlc)
	}
	return prices, nil
}
