package database_test

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/lunagic/database-go/database"
	"github.com/lunagic/database-go/database/internal/tester"
	"github.com/ory/dockertest/v3"
)

func getDockerDBAL(
	t *testing.T,
	defaultPort int,
	resourceGetter func(pool *dockertest.Pool) (*dockertest.Resource, error),
	driverGetter func(host string, port int) database.Driver,
) *database.DBAL {
	t.Helper()
	if testing.Short() {
		t.Skip("Skipping long-running test in short mode.")
	}

	var db *database.DBAL

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		t.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := resourceGetter(pool)
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	})

	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		dockerURL = "tcp://" + resource.GetHostPort(fmt.Sprintf("%d/tcp", defaultPort))
	}
	u, err := url.Parse(dockerURL)
	if err != nil {
		t.Fatalf("Error parsing docker URL: %s", err)
	}

	port := func() int {
		i, _ := strconv.Atoi(u.Port())
		return i
	}()

	if err := pool.Retry(func() error {
		var err error
		db, err = database.NewDBAL(driverGetter(u.Hostname(), port), database.WithLogger(tester.Logger(t)))
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to database: %s", err)
	}

	return db
}
