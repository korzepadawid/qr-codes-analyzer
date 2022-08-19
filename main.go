package main

import (
	"database/sql"
	"github.com/korzepadawid/qr-codes-analyzer/cache"
	"github.com/korzepadawid/qr-codes-analyzer/encode"
	"github.com/korzepadawid/qr-codes-analyzer/storage"
	"github.com/korzepadawid/qr-codes-analyzer/token"
	"github.com/korzepadawid/qr-codes-analyzer/util"
	"log"
	"time"

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

	server, err := api.NewServer(
		*cfg,
		db.NewStore(conn),
		token.NewJWTMaker("fsdfsdf", time.Hour),
		util.NewBCryptHasher(),
		storage.NewAWSS3FileStorageService(cfg),
		encode.NewQRCodeEncoder(),
		cache.NewRedisCache(),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = server.Run()

	if err != nil {
		log.Fatal(err)
	}
}
