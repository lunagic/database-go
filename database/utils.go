package database

import (
	"reflect"
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
