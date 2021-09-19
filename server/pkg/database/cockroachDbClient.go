package database

import (
	"context"
	"fmt"
	"log"
	"stonksio/pkg/common"
	"stonksio/pkg/config"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/google/uuid"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"

	"github.com/jackc/pgx/v4/pgxpool"
)

type CockroachDbClient struct {
	config config.DatabaseConfig
	db     *pgxpool.Pool
}

func NewCockroachDbClient(
	config *config.Config,
) *CockroachDbClient {
	// Connect to the stonksio database
	connConfig, err := pgxpool.ParseConfig(config.DatabaseConfig.ConnectionString)
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}

	conn, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	return &CockroachDbClient{
		config: config.DatabaseConfig,
		db:     conn,
	}
}

func (client *CockroachDbClient) closeConn() {
	client.db.Close()
}

func (client *CockroachDbClient) InsertPost(
	post common.Post,
) error {
	return crdbpgx.ExecuteTx(context.Background(), client.db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		//log.Printf("Creating post=%s\n", post)
		_, err := tx.Exec(context.Background(),
			`INSERT INTO post (id, username, "userpicurl", "body", "timestamp") VALUES ($1, $2, $3, $4, $5)`,
			uuid.New().String(), post.Username, post.UserPicUrl, post.Body, post.Timestamp)
		return err
	})
}

func (client *CockroachDbClient) InsertWallet(
	wallet common.Wallet,
) error {
	return crdbpgx.ExecuteTx(context.Background(), client.db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		log.Printf("inserting wallet. balance=%f", wallet.Balance)
		_, err := tx.Exec(context.Background(),
			`INSERT INTO wallet ("id", "username", "asset", "balance") VALUES ($1, $2, $3, $4)`,
			uuid.New().String(), wallet.Username, wallet.Asset, wallet.Balance)
		return err
	})
}

func (client *CockroachDbClient) deleteAllPosts() error {
	return crdbpgx.ExecuteTx(context.Background(), client.db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		log.Printf("Deleting all posts")
		_, err := tx.Exec(context.Background(), "DELETE FROM post")
		return err
	})
}

func (client *CockroachDbClient) deleteAllPrices() error {
	return crdbpgx.ExecuteTx(context.Background(), client.db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		log.Printf("Deleting all prices")
		_, err := tx.Exec(context.Background(), "DELETE FROM price")
		return err
	})
}

func (client *CockroachDbClient) GetPosts(n int) ([]common.Post, error) {
	rows, err := client.db.Query(context.Background(),
		`SELECT id, username, userpicurl, body, timestamp FROM post ORDER BY timestamp DESC LIMIT $1;`, n)
	if err != nil {
		return nil, fmt.Errorf("cannot query rows. err=%s", err)
	}
	posts := make([]common.Post, 0, n)
	defer rows.Close()
	for rows.Next() {
		post := common.Post{}
		if err := rows.Scan(&post.Id, &post.Username, &post.UserPicUrl, &post.Body, &post.Timestamp); err != nil {
			return nil, fmt.Errorf("cannot scan rows. err=%s", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (client *CockroachDbClient) GetBalance(
	asset, username string,
) (float32, error) {
	rows, err := client.db.Query(context.Background(),
		`SELECT balance FROM wallet WHERE username=$1 AND asset=$2;`, username, asset)
	if err != nil {
		return 0, fmt.Errorf("cannot query balance for username=%s. err=%s", username, err)
	}
	defer rows.Close()
	for rows.Next() {
		var balance float32
		if err := rows.Scan(&balance); err != nil {
			return 0, err
		}
		return balance, nil
	}
	//return 0, fmt.Errorf("no balance found for username=%s", username)
	return -1, nil
}

func (client *CockroachDbClient) GetPrices(
	asset string, window time.Duration,
) ([]common.Price, error) {
	if asset != "ETH" {
		return nil, fmt.Errorf("invalid asset=%s", asset)
	}
	startTime := time.Now().Add(-window)
	rows, err := client.db.Query(context.Background(),
		`SELECT tradePrice, timestamp FROM price WHERE asset=$1 AND timestamp > $2 ORDER BY timestamp;`, asset, startTime)
	if err != nil {
		return nil, fmt.Errorf("cannot query prices. err=%s", err)
	}
	prices := make([]common.Price, 0)
	defer rows.Close()
	for rows.Next() {
		price := common.Price{
			Asset: asset,
		}
		var tradePrice float32
		if err := rows.Scan(&tradePrice, &price.Timestamp); err != nil {
			return nil, err
		}
		price.TradePrice = tradePrice
		prices = append(prices, price)
	}
	return prices, nil
}

func (client *CockroachDbClient) GetOHLCs(count int) ([]common.OHLC, error) {
	rows, err := client.db.Query(context.Background(), `
		with intervals as (
		  select start, start + interval '1min' as end
		  from generate_series(
			date_trunc('minute', NOW()::timestamp - ($1-1) * interval '1 min'),
			date_trunc('minute', NOW()::timestamp),
			interval '1min')
	      as start
		)
		select distinct
		  intervals.start as date,
		  min(p.tradeprice) over w as low,
		  max(p.tradeprice) over w as high,
		  first_value(p.tradeprice) over w as open,
		  last_value(p.tradeprice) over w as close
		from
		  intervals
		  join price p on
		    p.timestamp >= intervals.start and
		    p.timestamp < intervals.end
		window w as (partition by intervals.start order by p.timestamp asc rows between unbounded preceding and unbounded following)
		order by intervals.start
	`, count)
	if err != nil {
		return nil, err
	}

	ohlcs := make([]common.OHLC, 0)
	defer rows.Close()
	for rows.Next() {
		ohlc := common.OHLC{}
		if err := rows.Scan(&ohlc.StartTime, &ohlc.Low, &ohlc.High, &ohlc.Open, &ohlc.Close); err != nil {
			return nil, err
		}
		ohlcs = append(ohlcs, ohlc)
	}
	return ohlcs, nil
}

func (client *CockroachDbClient) GetLatestPrice(
	asset string,
) (float32, error) {
	if asset != "ETH" {
		return 0, fmt.Errorf("invalid asset=%s", asset)
	}

	rows, err := client.db.Query(context.Background(),
		"SELECT tradePrice FROM price ORDER BY timestamp DESC LIMIT 1")
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		var tradePrice float32
		if err := rows.Scan(&tradePrice); err != nil {
			return 0, err
		}
		return tradePrice, nil
	}
	return 0, fmt.Errorf("no latest price found")
}

func (client *CockroachDbClient) InsertPrice(price common.Price) error {
	if price.Asset != "ETH" {
		return fmt.Errorf("invalid asset=%s", price.Asset)
	}
	return crdbpgx.ExecuteTx(context.Background(), client.db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		//log.Printf("Creating tradePrice=%f for asset=%s\n", tradePrice, asset)
		_, err := tx.Exec(context.Background(),
			"INSERT INTO price (id, asset, tradePrice, timestamp) VALUES ($1, $2, $3, $4)",
			uuid.New().String(), price.Asset, price.TradePrice, price.Timestamp)
		return err
	})
}
