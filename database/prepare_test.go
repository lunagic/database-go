package database_test

import "testing"

func Test_Prase_One(t *testing.T) {
	prepareTestHelper(t, PrepareTestCase{
		OriginalStatement: "SELECT * FROM user where id = :userID AND active = :active AND role IN (:roles) AND userID = :userID",
		OriginalParameters: map[string]any{
			":userID": 5,
			":active": 1,
			":roles": []string{
				"agent",
				"admin",
			},
		},
		ExpectedStatement: "SELECT * FROM user where id = ? AND active = ? AND role IN (?, ?) AND userID = ?",
		ExpectedParameters: []any{
			5,
			1,
			"agent",
			"admin",
			5,
		},
		ExpectedError: nil,
	})
}
