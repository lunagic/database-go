package database

import (
	"reflect"
	"regexp"
	"strings"
)

func loopOverStructFields(value reflect.Value, fieldHandler func(fieldDefinition reflect.StructField, fieldValue reflect.Value)) {
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldDefinition := value.Type().Field(i)

		if !fieldDefinition.IsExported() {
			continue
		}

		if fieldDefinition.Type.Kind() == reflect.Struct && fieldDefinition.Anonymous {
			loopOverStructFields(fieldValue, fieldHandler)

			continue
		}

		fieldHandler(fieldDefinition, fieldValue)
	}
}

type dbalTag struct {
	Column     string
	PrimaryKey bool
	ReadOnly   bool
}

func parseTag(tagString reflect.StructTag) dbalTag {
	parts := strings.Split(tagString.Get("db"), ",")

	tag := dbalTag{}

	for i, part := range parts {
		if i == 0 {
			tag.Column = part
			continue
		}

		if part == "readOnly" {
			tag.ReadOnly = true

			continue
		}

		if part == "primaryKey" {
			tag.PrimaryKey = true

			continue
		}
	}

	return tag
}

func prepare(statement string, parameters map[string]any) (string, []any, error) {
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
