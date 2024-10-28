package database_test

import (
	"testing"

	"github.com/lunagic/database-go/database"
	"github.com/ory/dockertest/v3"
)

func TestMySQL(t *testing.T) {
	runDriverTestSuite(
		t,
		getDockerDBAL(
			t,
			3306,
			func(pool *dockertest.Pool) (*dockertest.Resource, error) {
				return pool.Run("mariadb", "latest", []string{
					"MARIADB_ROOT_PASSWORD=root_password",
					"MARIADB_DATABASE=testing_database",
				})
			},
			func(host string, port int) database.Driver {
				return database.DriverMySQL{
					Username: "root",
					Name:     "testing_database",
					Password: "root_password",
					Hostname: host,
					Port:     port,
				}
			},
		),
	)
}
