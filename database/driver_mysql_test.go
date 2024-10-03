package database_test

import (
	"testing"

	"github.com/lunagic/database-go/database"
)

func Test_GenerateDelete_One(t *testing.T) {
	testHelperForQueries(t, TestCaseForQueryGeneration{
		Generator: database.MySQL{}.GenerateDelete,
		Input: User{
			ID: 1,
		},
		ExpectedStatement: "DELETE FROM `user` WHERE `id` = :id",
		ExpectedParameters: map[string]any{
			":id": UserID(1),
		},
	})
}

func Test_GenerateInsert_One(t *testing.T) {
	testHelperForQueries(t, TestCaseForQueryGeneration{
		Generator: database.MySQL{}.GenerateInsert,
		Input: User{
			ID:   1,
			Name: "Aaron",
		},
		ExpectedStatement: "INSERT INTO `user` (`id`, `name`) VALUES (:id, :name)",
		ExpectedParameters: map[string]any{
			":id":   UserID(1),
			":name": "Aaron",
		},
	})
}

func Test_GenerateUpdate_One(t *testing.T) {
	testHelperForQueries(t, TestCaseForQueryGeneration{
		Generator: database.MySQL{}.GenerateUpdate,
		Input: User{
			ID:   1,
			Name: "Aaron",
		},
		ExpectedStatement: "UPDATE `user` SET `name` = :name WHERE `id` = :id",
		ExpectedParameters: map[string]any{
			":id":   UserID(1),
			":name": "Aaron",
		},
	})
}

func Test_GenerateSelect_One(t *testing.T) {
	testHelperForQueries(t, TestCaseForQueryGeneration{
		Generator: func(entity database.Entity) (string, map[string]any, error) {
			query, err := database.MySQL{}.GenerateSelect(entity)
			if err != nil {
				return "", nil, err
			}

			return query.String(), map[string]any{}, nil
		},
		Input:              User{},
		ExpectedStatement:  "SELECT `id`, `name`, `updated_at` FROM `user`",
		ExpectedParameters: map[string]any{},
	})
}
