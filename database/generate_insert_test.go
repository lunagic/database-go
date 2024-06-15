package database_test

import (
	"testing"

	"github.com/lunagic/database-go/database"
)

func Test_GenerateInsert_One(t *testing.T) {
	testHelperForQueries(t, TestCaseForQueryGeneration{
		Generator: database.GenerateInsert,
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
