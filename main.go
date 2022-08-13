package main

import (
	"database/sql"
	"log"

	"github.com/korzepadawid/qr-codes-analyzer/api"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load("./")

	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(*cfg, store)

	if err != nil {
		log.Fatal(err)
	}

	err = server.Run()

	if err != nil {
		log.Fatal(err)
	}
}
