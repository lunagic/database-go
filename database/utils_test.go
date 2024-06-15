package database_test

import (
	"context"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"github.com/lunagic/database-go/database"
	"github.com/ory/dockertest/v3"
)

type PrepareTestCase struct {
	OriginalStatement  string
	OriginalParameters map[string]any
	ExpectedStatement  string
	ExpectedParameters []any
	ExpectedError      error
}

func prepareTestHelper(t *testing.T, testCase PrepareTestCase) {
	t.Helper()

	actualStatement, actualArgs, err := database.Prepare(testCase.OriginalStatement, testCase.OriginalParameters)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, testCase.ExpectedStatement, actualStatement) {
		return
	}

	if !assert.Equal(t, testCase.ExpectedParameters, actualArgs) {
		return
	}
}

func getDockerDBAL(t *testing.T) (context.Context, *database.DBAL) {
	var db *database.DBAL

	config := database.Config{
		Username: "root",
		Password: "secret",
		Name:     "testingDB",
	}

	mysql.SetLogger(log.New(io.Discard, "", log.LstdFlags))

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
	resource, err := pool.Run("mariadb", "latest", []string{
		"MARIADB_ROOT_PASSWORD=" + config.Password,
		"MARIADB_DATABASE=" + config.Name,
	})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		dockerURL = "tcp://" + resource.GetHostPort("3306/tcp")
	}
	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}

	config.Hostname = u.Hostname()
	config.Port = func() int {
		i, _ := strconv.Atoi(u.Port())
		return i
	}()

	if err := pool.Retry(func() error {
		var err error
		db, err = database.NewDBAL(config, log.New(os.Stdout, "", log.LstdFlags))
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to database: %s", err)
	}

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	})

	t.Log("Starting test")

	return context.Background(), db
}

type TestCaseForQueryGeneration struct {
	Input              database.Entity
	ExpectedStatement  string
	ExpectedParameters map[string]any
	ExpectedError      error
	Generator          func(entity database.Entity) (string, map[string]any, error)
}

func testHelperForQueries(t *testing.T, testCase TestCaseForQueryGeneration) {
	t.Helper()

	actualStatement, actualParameters, actualErr := testCase.Generator(testCase.Input)
	if !assert.Equal(t, testCase.ExpectedError, actualErr) {
		return
	}

	if !assert.Equal(t, testCase.ExpectedStatement, actualStatement) {
		return
	}

	if !assert.Equal(t, testCase.ExpectedParameters, actualParameters) {
		return
	}
}
