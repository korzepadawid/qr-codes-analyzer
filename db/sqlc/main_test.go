package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/korzepadawid/qr-codes-analyzer/config"
	_ "github.com/lib/pq"
)

var testDB *sql.DB
var testQueries *Queries

func TestMain(m *testing.M) {
	var err error
	cfg, err := config.Load("../..")

	if err != nil {
		log.Fatal(err)
	}

	testDB, err = sql.Open("postgres", cfg.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
