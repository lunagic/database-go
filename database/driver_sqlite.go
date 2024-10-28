package database

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DriverSQLite struct {
	Path string
}

func (driver DriverSQLite) DSN() string {
	return fmt.Sprintf("file:%s?cache=shared", driver.Path)
}

func (driver DriverSQLite) Driver() string {
	return "sqlite3"
}

func (driver DriverSQLite) ShowCreateTable(tableName string) string {
	return fmt.Sprintf("SELECT '' as `Table`, sql as `Create Table` FROM sqlite_schema WHERE name='%s';", tableName)
}

func (driver DriverSQLite) translateType(column TableColumn) string {
	switch column.Type.Kind() {
	case reflect.String:
		return "text NOT NULL"
	case reflect.Uint64:
		return "integer NOT NULL"
	}

	if column.Type.String() == "time.Time" {
		return "datetime NOT NULL DEFAULT CURRENT_TIMESTAMP"
	}

	if column.Type.String() == "*time.Time" {
		return "datetime NULL DEFAULT NULL"
	}

	log.Fatalf("unsupported type: %s (%s)", column.Type.String(), column.Type.Kind().String())
	return ""
}

func (driver DriverSQLite) CreateTable(t Table) string {
	statements := []string{}

	for _, column := range t.Columns {
		statement := fmt.Sprintf(
			"`%s` %s",
			column.Name,
			driver.translateType(column),
		)

		if column.PrimaryKey {
			statement = fmt.Sprintf("%s PRIMARY KEY", statement)
		}

		statements = append(statements, statement)
	}

	return fmt.Sprintf("CREATE TABLE `%s` (%s)", t.Name, strings.Join(statements, ", "))

}

func (driver DriverSQLite) Insert(entity Entity) (string, map[string]any, error) {
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

		if tag.AutoIncrement {
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

func (driver DriverSQLite) Delete(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}

func (driver DriverSQLite) Save(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}

func (driver DriverSQLite) Select(entity Entity) (Query, error) {
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

func (driver DriverSQLite) Update(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}

func (driver DriverSQLite) TableFromEntity(e Entity) Table {
	return Table{}
}

func (driver DriverSQLite) ParseCreateTable(s string) Table {
	return Table{}
}

func (driver DriverSQLite) RenameTable() string {
	return ""
}

func (driver DriverSQLite) DropTable(t Table) string {
	return ""
}

func (driver DriverSQLite) AddColumn(c TableColumn) string {
	return ""
}

func (driver DriverSQLite) AlterColumn(c TableColumn) string {
	return ""
}

func (driver DriverSQLite) DropColumn(c TableColumn) string {
	return ""
}

func (driver DriverSQLite) AddIndex(c TableIndex) string {
	return ""
}

func (driver DriverSQLite) DropIndex(c TableIndex) string {
	return ""
}

func (driver DriverSQLite) AddForeignKey(c TableIndex) string {
	return ""
}

func (driver DriverSQLite) DropForeignKey(c TableIndex) string {
	return ""
}
