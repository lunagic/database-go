package database_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/lunagic/database-go/database"
)

func TestSQLite(t *testing.T) {
	dbal, err := database.NewDBAL(
		database.DriverSQLite{
			Path: fmt.Sprintf("%s/database.sqlite", t.TempDir()),
		},
		database.WithLogger(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		t.Fatal(err)
	}
	runDriverTestSuite(t, dbal)
}
