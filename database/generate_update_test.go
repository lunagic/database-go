package database_test

import (
	"testing"

	"github.com/lunagic/database-go/database"
)

func Test_GenerateUpdate_One(t *testing.T) {
	testHelperForQueries(t, TestCaseForQueryGeneration{
		Generator: database.GenerateUpdate,
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
