package database_test

// type TestCaseForQueryGeneration struct {
// 	Input              database.Entity
// 	ExpectedStatement  string
// 	ExpectedParameters map[string]any
// 	ExpectedError      error
// 	Generator          func(entity database.Entity) (string, map[string]any, error)
// }

// func testHelperForQueries(t *testing.T, testCase TestCaseForQueryGeneration) {
// 	t.Helper()

// 	actualStatement, actualParameters, actualErr := testCase.Generator(testCase.Input)
// 	if !assert.Equal(t, testCase.ExpectedError, actualErr) {
// 		return
// 	}

// 	if !assert.Equal(t, testCase.ExpectedStatement, actualStatement) {
// 		return
// 	}

// 	if !assert.Equal(t, testCase.ExpectedParameters, actualParameters) {
// 		return
// 	}
// }

// type PrepareTestCase struct {
// 	OriginalStatement  string
// 	OriginalParameters map[string]any
// 	ExpectedStatement  string
// 	ExpectedParameters []any
// 	ExpectedError      error
// }

// func prepareTestHelper(t *testing.T, testCase PrepareTestCase) {
// 	t.Helper()

// 	actualStatement, actualArgs, err := database.Prepare(testCase.OriginalStatement, testCase.OriginalParameters)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !assert.Equal(t, testCase.ExpectedStatement, actualStatement) {
// 		return
// 	}

// 	if !assert.Equal(t, testCase.ExpectedParameters, actualArgs) {
// 		return
// 	}
// }
// func Test_Prase_One(t *testing.T) {
// 	prepareTestHelper(t, PrepareTestCase{
// 		OriginalStatement: "SELECT * FROM user where id = :userID AND active = :active AND role IN (:roles) AND userID = :userID",
// 		OriginalParameters: map[string]any{
// 			":userID": 5,
// 			":active": 1,
// 			":roles": []string{
// 				"agent",
// 				"admin",
// 			},
// 		},
// 		ExpectedStatement: "SELECT * FROM user where id = ? AND active = ? AND role IN (?, ?) AND userID = ?",
// 		ExpectedParameters: []any{
// 			5,
// 			1,
// 			"agent",
// 			"admin",
// 			5,
// 		},
// 		ExpectedError: nil,
// 	})
// }
