package database_test

import (
	"testing"

	"github.com/lunagic/database-go/database"
	"github.com/ory/dockertest/v3"
)

func Test_DriverPostgres_17(t *testing.T) {
	t.Parallel()
	testDriverPostgres(t, "17")
}

func Test_DriverPostgres_16(t *testing.T) {
	t.Parallel()
	testDriverPostgres(t, "16")
}

func Test_DriverPostgres_15(t *testing.T) {
	t.Parallel()
	testDriverPostgres(t, "15")
}

func testDriverPostgres(t *testing.T, version string) {
	t.SkipNow() // TODO: remove this for once driver is implemented
	runDriverTestSuite(
		t,
		getDockerDBAL(
			t,
			5432,
			func(pool *dockertest.Pool) (*dockertest.Resource, error) {
				return pool.Run("postgres", version, []string{
					"POSTGRES_USER=root_user",
					"POSTGRES_PASSWORD=root_password",
					"POSTGRES_DB=testing_database",
				})
			},
			func(host string, port int) database.Driver {
				return database.DriverPostgres{
					Hostname: host,
					Port:     port,
					Username: "root_user",
					Password: "root_password",
					Name:     "testing_database",
				}
			},
		),
	)
}
