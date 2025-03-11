package database_test

import (
	"testing"

	"github.com/lunagic/database-go/database"
	"github.com/ory/dockertest/v3"
)

func Test_DriverMySQL_8(t *testing.T) {
	t.Skip()
	t.Parallel()
	testDriverMySQL(t, "mysql", "8")
}

func Test_DriverMySQL_MariaDB_11_4(t *testing.T) {
	t.Parallel()
	testDriverMySQL(t, "mariadb", "11.4")
}

func Test_DriverMySQL_MariaDB_10_11(t *testing.T) {
	t.Parallel()
	testDriverMySQL(t, "mariadb", "10.11")
}

func Test_DriverMySQL_MariaDB_10_6(t *testing.T) {
	t.Parallel()
	testDriverMySQL(t, "mariadb", "10.6")
}

func testDriverMySQL(t *testing.T, server string, version string) {
	t.Helper()
	runDriverTestSuite(
		t,
		getDockerDBAL(
			t,
			3306,
			func(pool *dockertest.Pool) (*dockertest.Resource, error) {
				return pool.Run(server, version, []string{
					"MYSQL_ROOT_PASSWORD=root_password",
					"MYSQL_DATABASE=testing_database",
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
