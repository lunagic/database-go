package database

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type DriverMySQL struct {
	Hostname string
	Port     int
	Username string
	Password string
	Name     string
}

func (driver DriverMySQL) Driver() string {
	_ = mysql.SetLogger(log.New(io.Discard, "", log.LstdFlags))

	return "mysql"
}

func (driver DriverMySQL) DSN() string {
	return fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?parseTime=true",
		driver.Username,
		driver.Password,
		driver.Hostname,
		driver.Port,
		driver.Name,
	)
}

func (driver DriverMySQL) Insert(entity Entity) (string, map[string]any, error) {
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

func (driver DriverMySQL) Delete(entity Entity) (string, map[string]any, error) {
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

func (driver DriverMySQL) Save(entity Entity) (string, map[string]any, error) {
	return "", map[string]any{}, nil
}

func (driver DriverMySQL) Select(entity Entity) (Query, error) {
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

func (driver DriverMySQL) Update(entity Entity) (string, map[string]any, error) {
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

func (driver DriverMySQL) TableFromEntity(e Entity) Table {
	return Table{}
}

func (driver DriverMySQL) ParseCreateTable(s string) Table {
	return Table{}
}

func (driver DriverMySQL) CreateTable(t Table) string {
	return ""
}

func (driver DriverMySQL) RenameTable() string {
	return ""
}

func (driver DriverMySQL) DropTable(t Table) string {
	return ""
}

func (driver DriverMySQL) AddColumn(c TableColumn) string {
	return ""
}

func (driver DriverMySQL) AlterColumn(c TableColumn) string {
	return ""
}

func (driver DriverMySQL) DropColumn(c TableColumn) string {
	return ""
}

func (driver DriverMySQL) AddIndex(c TableIndex) string {
	return ""
}

func (driver DriverMySQL) DropIndex(c TableIndex) string {
	return ""
}

func (driver DriverMySQL) AddForeignKey(c TableIndex) string {
	return ""
}

func (driver DriverMySQL) DropForeignKey(c TableIndex) string {
	return ""
}
