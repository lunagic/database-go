package database

import (
	"reflect"
	"regexp"
	"strings"
)

func Prepare(statement string, parameters map[string]any) (string, []any, error) {
	paramFinder := regexp.MustCompile(`(?m):\w+`)
	spaceFinder := regexp.MustCompile(`(?m)\s^\s+`)

	statement = spaceFinder.ReplaceAllString(statement, " ")

	args := []any{}

	newStatement := paramFinder.ReplaceAllStringFunc(statement, func(s string) string {
		parameterValue, found := parameters[s]
		if !found {
			return s
		}

		rt := reflect.TypeOf(parameterValue)
		if rt.Kind() == reflect.Array || rt.Kind() == reflect.Slice {
			localArgs := []string{}

			valueOf := reflect.ValueOf(parameterValue)
			for i := 0; i < valueOf.Len(); i++ {
				localArgs = append(localArgs, "?")
				args = append(args, valueOf.Index(i).Interface())
			}

			return strings.Join(localArgs, ", ")
		}

		args = append(args, parameterValue)

		return "?"
	})

	return newStatement, args, nil
}
