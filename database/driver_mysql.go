package database

import (
	"fmt"
	"reflect"
	"strings"
)

type MySQL struct {
	Hostname string
	Port     int
	Username string
	Password string
	Name     string
}

func (g MySQL) Driver() string {
	return "mysql"
}

func (q MySQL) DSN() string {
	return fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?parseTime=true",
		q.Username,
		q.Password,
		q.Hostname,
		q.Port,
		q.Name,
	)
}

func (g MySQL) GenerateInsert(entity Entity) (string, map[string]any, error) {
	columns := []string{}
	values := []string{}
	parameters := map[string]any{}

	loopOverStructFields(reflect.ValueOf(entity), func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
		tag := parseTag(fieldDefinition.Tag)
		if tag.Column == "" {
			return
		}

		if tag.ReadOnly {
			return
		}

		column := fmt.Sprintf("`%s`", tag.Column)
		value := fmt.Sprintf(":%s", tag.Column)

		columns = append(columns, column)
		values = append(values, value)

		parameters[value] = fieldValue.Interface()

	})

	return fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s)",
		entity.EntityInformation().TableName,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	), parameters, nil
}

func (g MySQL) GenerateDelete(entity Entity) (string, map[string]any, error) {
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

func (g MySQL) GenerateSave(entity Entity) (string, map[string]any, error) {
	return "", map[string]any{}, nil
}

func (g MySQL) GenerateSelect(entity Entity) (Query, error) {
	selects := []string{}

	loopOverStructFields(reflect.ValueOf(entity), func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
		tag := parseTag(fieldDefinition.Tag)
		if tag.Column == "" {
			return
		}

		column := fmt.Sprintf("`%s`", tag.Column)
		selects = append(selects, column)
	})

	return Query{
		Select: selects,
		From:   entity.EntityInformation().TableName,
	}, nil
}

func (g MySQL) GenerateUpdate(entity Entity) (string, map[string]any, error) {
	updates := []string{}
	primaryKeys := []string{}

	parameters := map[string]any{}

	loopOverStructFields(reflect.ValueOf(entity), func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
		tag := parseTag(fieldDefinition.Tag)
		if tag.Column == "" {
			return
		}

		if tag.ReadOnly {
			return
		}

		column := fmt.Sprintf("`%s`", tag.Column)
		value := fmt.Sprintf(":%s", tag.Column)

		condition := fmt.Sprintf("%s = %s", column, value)

		parameters[value] = fieldValue.Interface()

		if tag.PrimaryKey {
			primaryKeys = append(primaryKeys, condition)
		} else {
			updates = append(updates, condition)
		}

	})

	return fmt.Sprintf(
		"UPDATE `%s` SET %s WHERE %s",
		entity.EntityInformation().TableName,
		strings.Join(updates, ", "),
		strings.Join(primaryKeys, " AND "),
	), parameters, nil
}
