package database_test

import (
	"testing"

	"github.com/lunagic/database-go/database"
)

func Test_GenerateSelect_One(t *testing.T) {
	testHelperForQueries(t, TestCaseForQueryGeneration{
		Generator: func(entity database.Entity) (string, map[string]any, error) {
			query, err := database.GenerateSelect(entity)
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
