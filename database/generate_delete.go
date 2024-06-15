package database

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateDelete(entity Entity) (string, map[string]any, error) {
	primaryKeys := []string{}

	parameters := map[string]any{}

	loopOverStructFields(reflect.ValueOf(entity), func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
		tag := parseTag(fieldDefinition.Tag)
		if tag.Column == "" {
			return
		}

		if !tag.PrimaryKey {
			return
		}

		column := fmt.Sprintf("`%s`", tag.Column)
		value := fmt.Sprintf(":%s", tag.Column)

		primaryKeys = append(primaryKeys, fmt.Sprintf("%s = %s", column, value))
		parameters[value] = fieldValue.Interface()
	})

	return fmt.Sprintf(
		"DELETE FROM `%s` WHERE %s",
		entity.EntityInformation().TableName,
		strings.Join(primaryKeys, " AND "),
	), parameters, nil
}
