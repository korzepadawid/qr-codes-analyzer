package db

import (
	"database/sql"
	"github.com/korzepadawid/qr-codes-analyzer/config"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testDB *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
	var err error
	cfg, err := config.Load("../..")

	if err != nil {
		log.Fatal(err)
	}

	testDB, err = sql.Open(cfg.DBDriver, cfg.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
