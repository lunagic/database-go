package database_test

import (
	"testing"

	"github.com/lunagic/database-go/database"
	"github.com/ory/dockertest/v3"
)

func Test_DriverMySQL_11_4(t *testing.T) {
	t.Parallel()
	testDriverMySQL(t, "11.4")
	t.Fail()
}

func Test_DriverMySQL_10_11(t *testing.T) {
	t.Parallel()
	testDriverMySQL(t, "10.11")
}

func Test_DriverMySQL_10_6(t *testing.T) {
	t.Parallel()
	testDriverMySQL(t, "10.6")
}

func testDriverMySQL(t *testing.T, version string) {
	t.Helper()
	runDriverTestSuite(
		t,
		getDockerDBAL(
			t,
			3306,
			func(pool *dockertest.Pool) (*dockertest.Resource, error) {
				return pool.Run("mariadb", version, []string{
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
