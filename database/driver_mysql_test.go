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

// func Test_GenerateDelete_One(t *testing.T) {
// 	testHelperForQueries(t, TestCaseForQueryGeneration{
// 		Generator: database.MySQL{}.Delete,
// 		Input: User{
// 			ID: 1,
// 		},
// 		ExpectedStatement: "DELETE FROM `user` WHERE `id` = :id",
// 		ExpectedParameters: map[string]any{
// 			":id": UserID(1),
// 		},
// 	})
// }

// func Test_GenerateInsert_One(t *testing.T) {
// 	testHelperForQueries(t, TestCaseForQueryGeneration{
// 		Generator: database.MySQL{}.Insert,
// 		Input: User{
// 			ID:   1,
// 			Name: "Aaron",
// 		},
// 		ExpectedStatement: "INSERT INTO `user` (`id`, `name`) VALUES (:id, :name)",
// 		ExpectedParameters: map[string]any{
// 			":id":   UserID(1),
// 			":name": "Aaron",
// 		},
// 	})
// }

// func Test_GenerateUpdate_One(t *testing.T) {
// 	testHelperForQueries(t, TestCaseForQueryGeneration{
// 		Generator: database.MySQL{}.Update,
// 		Input: User{
// 			ID:   1,
// 			Name: "Aaron",
// 		},
// 		ExpectedStatement: "UPDATE `user` SET `name` = :name WHERE `id` = :id",
// 		ExpectedParameters: map[string]any{
// 			":id":   UserID(1),
// 			":name": "Aaron",
// 		},
// 	})
// }

// func Test_GenerateSelect_One(t *testing.T) {
// 	testHelperForQueries(t, TestCaseForQueryGeneration{
// 		Generator: func(entity database.Entity) (string, map[string]any, error) {
// 			query, err := database.MySQL{}.Select(entity)
// 			if err != nil {
// 				return "", nil, err
// 			}

// 			return query.String(), map[string]any{}, nil
// 		},
// 		Input:              User{},
// 		ExpectedStatement:  "SELECT `id`, `name`, `updated_at` FROM `user`",
// 		ExpectedParameters: map[string]any{},
// 	})
// }
