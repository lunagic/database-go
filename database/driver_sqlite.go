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

func (driver DriverSQLite) Driver() string {
	return "sqlite3"
}

func (driver DriverSQLite) DSN() string {
	return fmt.Sprintf("file:%s?cache=shared", driver.Path)
}

func (driver DriverSQLite) ShowCreateTable(tableName string) string {
	return fmt.Sprintf("SELECT sql as `Create Table` FROM sqlite_schema WHERE name='%s';", tableName)
}

func (driver DriverSQLite) ParseCreateTable(s string) Table {
	columns := []TableColumn{}

	for _, columnString := range strings.Split(strings.TrimSuffix(strings.SplitN(s, "(", 2)[1], ")"), ", ") {
		columns = append(columns, driver.stringToColumn(columnString))
	}

	return Table{
		Name:    "", // TODO: make this work
		Columns: columns,
	}
}

func (driver DriverSQLite) TableFromEntity(entity Entity) Table {
	table := Table{
		Name: entity.EntityInformation().TableName,
	}

	loopOverStructFields(reflect.ValueOf(entity), func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
		tag := parseTag(fieldDefinition.Tag)
		if tag.Column == "" {
			return
		}

		column := TableColumn{
			Name:          tag.Column,
			AutoIncrement: tag.AutoIncrement,
			Type:          driver.translateType(fieldDefinition.Type),
			PrimaryKey:    tag.PrimaryKey,
			HasDefault:    tag.HasDefault,
			Default:       tag.Default,
			Nullable:      reflect.Pointer == fieldValue.Kind(),
		}

		table.Columns = append(table.Columns, column)
	})

	return table
}

func (driver DriverSQLite) CreateTable(t Table) string {
	statements := []string{}

	for _, column := range t.Columns {
		statement := driver.columnToString(column)

		statements = append(statements, statement)
	}

	return fmt.Sprintf("CREATE TABLE `%s` (%s)", t.Name, strings.Join(statements, ", "))
}

func (driver DriverSQLite) RenameTable(t Table) string {
	return ""
}

func (driver DriverSQLite) DropTable(t Table) string {
	return fmt.Sprintf("DROP TABLE `%s`;", t.Name)
}

func (driver DriverSQLite) AddColumn(t Table, c TableColumn) string {
	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", t.Name, driver.columnToString(c))
}

func (driver DriverSQLite) DropColumn(t Table, c TableColumn) string {
	return fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", t.Name, c.Name)
}

func (driver DriverSQLite) AlterColumn(t Table, c TableColumn) string {
	return ""
}

func (driver DriverSQLite) AddIndex(t Table, c TableIndex) string {
	return ""
}

func (driver DriverSQLite) DropIndex(t Table, c TableIndex) string {
	return ""
}

func (driver DriverSQLite) AddForeignKey(t Table, c TableForeignKey) string {
	return ""
}

func (driver DriverSQLite) DropForeignKey(t Table, c TableForeignKey) string {
	return ""
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

func (driver DriverSQLite) Save(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}

func (driver DriverSQLite) Delete(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}

func (driver DriverSQLite) stringToColumn(columnString string) TableColumn {
	return TableColumn{
		PrimaryKey:    strings.Contains(columnString, "PRIMARY KEY"),
		AutoIncrement: strings.Contains(columnString, "PRIMARY KEY") && strings.Contains(columnString, "integer"),
		Nullable:      !strings.Contains(columnString, "NOT NULL"),
		Name:          strings.TrimPrefix(strings.Split(columnString, "` ")[0], "`"),
	}
}

func (driver DriverSQLite) columnToString(t TableColumn) string {
	column := fmt.Sprintf(
		"`%s` %s",
		t.Name,
		t.Type,
	)

	if t.Nullable {
		column = fmt.Sprintf("%s NULL", column)
	} else {
		column = fmt.Sprintf("%s NOT NULL", column)
	}

	if t.HasDefault {
		switch t.Default {
		case "NULL":
			column = fmt.Sprintf("%s DEFAULT NULL", column)
		case "CURRENT_TIMESTAMP":
			column = fmt.Sprintf("%s DEFAULT CURRENT_TIMESTAMP", column)
		default:
			column = fmt.Sprintf("%s DEFAULT '%s'", column, t.Default) // TODO: need to make sure we escape single quotes
		}
	} else if t.Type == "text" && !t.Nullable {
		column = fmt.Sprintf("%s DEFAULT ''", column)
	}

	if t.PrimaryKey {
		column = fmt.Sprintf("%s PRIMARY KEY", column)
	}

	return column
}

func (driver DriverSQLite) translateType(t reflect.Type) string {
	// Read through the pointer/slice/map
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "text"
	case reflect.Bool:
		return "integer"
	case reflect.Uint64:
		return "integer"
	}

	if t.String() == "time.Time" {
		return "datetime"
	}

	log.Fatalf("unsupported type: %s (%s)", t.String(), t.Kind().String())
	return ""
}
