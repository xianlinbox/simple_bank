package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	dbSource = "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	
	pgConn, err := pgx.Connect(context.Background(), dbSource)

	if (err != nil) {
		log.Fatal("cannot connect to db: ", err)
	}
	testQueries = New(pgConn)
	os.Exit(m.Run())
}