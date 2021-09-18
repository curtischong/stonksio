package main

import (
	"context"
	"log"

	"stonksio/pkg/config"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	yaml "gopkg.in/yaml.v3"
)

func insertRows(ctx context.Context, tx pgx.Tx, accts [4]uuid.UUID) error {
	print(yaml.AliasNode)
	// Insert four rows into the "post" table.
	log.Println("Creating new rows...")
	if _, err := tx.Exec(ctx,
		"INSERT INTO post (id, message) VALUES ($1, $2), ($3, $4), ($5, $6), ($7, $8)", accts[0], "250", accts[1], "100", accts[2], "500", accts[3], "300"); err != nil {
		return err
	}
	return nil
}

func printPosts(conn *pgx.Conn) error {
	rows, err := conn.Query(context.Background(), "SELECT id, message FROM post")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var message string
		if err := rows.Scan(&id, &message); err != nil {
			log.Fatal(err)
			return err
		}
		log.Printf("%s: %d\n", id, message)
	}
	return nil
}

func deleteRows(ctx context.Context, tx pgx.Tx, one uuid.UUID, two uuid.UUID) error {
	// Delete two rows into the "post" table.
	log.Printf("Deleting rows with IDs %s and %s...", one, two)
	if _, err := tx.Exec(ctx, "DELETE FROM post WHERE id IN ($1, $2)", one, two); err != nil {
		return err
	}
	log.Println("Deleted rows")
	return nil
}

func deleteAll(ctx context.Context, tx pgx.Tx) error {
	// Delete two rows into the "post" table.
	log.Printf("Deleting rows")
	if _, err := tx.Exec(ctx, "DELETE FROM post"); err != nil {
		return err
	}
	log.Println("Deleted all rows")
	return nil
}

const configPath = "./config.yaml"

func main() {
	stonksConfig, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("couldn't load config path=%s, err=%s", configPath, err)
	}

	// Connect to the stonksio database
	connConfig, err := pgx.ParseConfig(stonksConfig.DatabaseConfig.ConnectionString)
	connConfig.Database = "stonksio"
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer conn.Close(context.Background())

	// Insert initial rows
	var posts [4]uuid.UUID
	for i := 0; i < len(posts); i++ {
		posts[i] = uuid.New()
	}

	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return insertRows(context.Background(), tx, posts)
	})
	if err != nil {
		log.Fatal("error: ", err)
	}
	log.Println("New rows created.")

	// Print out the balances
	log.Println("Initial balances:")

	err = printPosts(conn)
	if err != nil {
		log.Println("error: ", err)
	}

	// Delete rows
	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return deleteRows(context.Background(), tx, posts[0], posts[1])
	})
	if err != nil {
		log.Fatal("error: ", err)
	}
	err = printPosts(conn)
	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return deleteAll(context.Background(), tx)
	})
}
