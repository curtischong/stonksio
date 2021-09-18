package main

import (
	"context"
	"log"
	"net/http"
	"stonksio/pkg/request"

	"stonksio/pkg/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

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
	config, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("couldn't load config path=%s, err=%s", configPath, err)
	}

	requestHandler := request.NewRequestHandler(config)
	http.HandleFunc("/get/ohlc/eth", requestHandler.HandleGetOhlc)
}
