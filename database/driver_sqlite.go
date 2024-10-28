package database

import "fmt"

type DriverSQLite struct {
	Path string
}

func (driver DriverSQLite) DSN() string {
	return fmt.Sprintf("file:%s?cache=shared", driver.Path)
}

func (driver DriverSQLite) Driver() string {
	return "sqlite3"
}

func (driver DriverSQLite) Insert(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}

func (driver DriverSQLite) Delete(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}

func (driver DriverSQLite) Save(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}

func (driver DriverSQLite) Select(entity Entity) (Query, error) {
	return Query{}, nil
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

func (driver DriverSQLite) CreateTable(t Table) string {
	return ""
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
