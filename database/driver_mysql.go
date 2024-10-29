package database

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"regexp"
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

func (driver DriverMySQL) ShowCreateTable(tableName string) string {
	return fmt.Sprintf("SHOW CREATE TABLE `%s`;", tableName)
}

func (driver DriverMySQL) ParseCreateTable(s string) Table {
	re := regexp.MustCompile(`(?m)\n\s+`)
	s = re.ReplaceAllString(s, " ")
	columns := []TableColumn{}

	for _, columnString := range strings.Split(strings.TrimSuffix(strings.SplitN(s, "(", 2)[1], ")"), ", ") {
		columnString = strings.TrimSpace(columnString)
		if !strings.HasPrefix(columnString, "`") {
			continue
		}

		columns = append(columns, driver.stringToColumn(columnString))
	}

	return Table{
		Name:    "", // TODO: make this work
		Columns: columns,
	}
}

func (driver DriverMySQL) TableFromEntity(e Entity) Table {
	table := Table{
		Name: e.EntityInformation().TableName,
	}

	loopOverStructFields(reflect.ValueOf(e), func(fieldDefinition reflect.StructField, fieldValue reflect.Value) {
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

func (driver DriverMySQL) CreateTable(t Table) string {
	// primaryKeyName := ""

	statements := []string{}
	for _, column := range t.Columns {
		// if column.PrimaryKey {
		// 	// primaryKeyName = column.Name
		// }

		statements = append(statements, driver.columnToString(column))
	}

	// if primaryKeyName != "" {
	// 	statements = append(statements, fmt.Sprintf(
	// 		"PRIMARY KEY (`%s`)",
	// 		primaryKeyName,
	// 	))
	// }

	return fmt.Sprintf("CREATE TABLE `%s` (%s)", t.Name, strings.Join(statements, ", "))
}

func (driver DriverMySQL) RenameTable(t Table) string {
	return ""
}

func (driver DriverMySQL) DropTable(t Table) string {
	return fmt.Sprintf("DROP TABLE `%s`;", t.Name)
}

func (driver DriverMySQL) AddColumn(t Table, c TableColumn) string {
	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", t.Name, driver.columnToString(c))
}

func (driver DriverMySQL) DropColumn(t Table, c TableColumn) string {
	return fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", t.Name, c.Name)
}

func (driver DriverMySQL) AlterColumn(t Table, c TableColumn) string {
	return ""
}

func (driver DriverMySQL) AddIndex(t Table, c TableIndex) string {
	return ""
}

func (driver DriverMySQL) DropIndex(t Table, c TableIndex) string {
	return ""
}

func (driver DriverMySQL) AddForeignKey(t Table, c TableForeignKey) string {
	return ""
}

func (driver DriverMySQL) DropForeignKey(t Table, c TableForeignKey) string {
	return ""
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

func (driver DriverMySQL) Save(entity Entity) (string, map[string]any, error) {
	return "", map[string]any{}, nil
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

func (driver DriverMySQL) stringToColumn(columnString string) TableColumn {
	return TableColumn{
		PrimaryKey:    strings.Contains(columnString, "PRIMARY KEY"),
		AutoIncrement: strings.Contains(columnString, "PRIMARY KEY") && strings.Contains(columnString, "integer"),
		Nullable:      !strings.Contains(columnString, "NOT NULL"),
		Name:          strings.TrimPrefix(strings.Split(columnString, "` ")[0], "`"),
	}
}

func (driver DriverMySQL) columnToString(t TableColumn) string {
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
		if t.Default == "NULL" {
			column = fmt.Sprintf("%s DEFAULT NULL", column)
		} else if t.Default == "CURRENT_TIMESTAMP" {
			column = fmt.Sprintf("%s DEFAULT CURRENT_TIMESTAMP", column)
		} else {
			column = fmt.Sprintf("%s DEFAULT '%s'", column, t.Default) // TODO: need to make sure we escape single quotes
		}
	} else if t.Type == "text" && t.Nullable == false {
		column = fmt.Sprintf("%s DEFAULT ''", column)
	}

	if t.PrimaryKey {
		column = fmt.Sprintf("%s PRIMARY KEY", column)
	}
	if t.AutoIncrement {
		column = fmt.Sprintf("%s AUTO_INCREMENT", column)
	}

	return column
}

func (driver DriverMySQL) translateType(t reflect.Type) string {
	// Read through the pointer/slice/map
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "text"
	case reflect.Uint64:
		return "int(20) unsigned"
	}

	if t.String() == "time.Time" {
		return "datetime"
	}

	log.Fatalf("unsupported type: %s (%s)", t.String(), t.Kind().String())
	return ""
}
