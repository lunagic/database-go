package database_test

import (
	"testing"

	"github.com/lunagic/database-go/database"
)

func Test_GenerateDelete_One(t *testing.T) {
	testHelperForQueries(t, TestCaseForQueryGeneration{
		Generator: database.GenerateDelete,
		Input: User{
			ID: 1,
		},
		ExpectedStatement: "DELETE FROM `user` WHERE `id` = :id",
		ExpectedParameters: map[string]any{
			":id": UserID(1),
		},
	})
}
