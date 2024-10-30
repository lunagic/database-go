package database_test

import (
	"fmt"
	"testing"

	"github.com/lunagic/database-go/database"
	"github.com/lunagic/database-go/database/internal/tester"
)

func Test_DriverSQLite(t *testing.T) {
	dbal, err := database.NewDBAL(
		database.DriverSQLite{
			Path: fmt.Sprintf("%s/database.sqlite", t.TempDir()),
		},
		database.WithLogger(tester.Logger(t)),
	)
	if err != nil {
		t.Fatal(err)
	}
	runDriverTestSuite(t, dbal)
}
