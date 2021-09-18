package main

import (
	"context"
	"log"
	"stonksio/pkg/database"

	"stonksio/pkg/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

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

const configPath = "./config.yaml"

func main() {
	stonksConfig, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("couldn't load config path=%s, err=%s", configPath, err)
	}

	cockroachDbClient := database.NewCockroachDbClient(stonksConfig)

}
