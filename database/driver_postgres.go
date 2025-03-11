package database

import (
	"fmt"

	_ "github.com/lib/pq"
)

type DriverPostgres struct {
	Hostname string
	Port     int
	Username string
	Password string
	Name     string
}

func (driver DriverPostgres) Driver() string {
	return "postgres"
}

func (driver DriverPostgres) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		driver.Hostname,
		driver.Port,
		driver.Username,
		driver.Password,
		driver.Name,
	)
}

func (driver DriverPostgres) ShowCreateTable(tableName string) string {
	return ""
}

func (driver DriverPostgres) ParseCreateTable(s string) Table {
	return Table{}
}

func (driver DriverPostgres) TableFromEntity(entity Entity) Table {
	return Table{}
}

func (driver DriverPostgres) CreateTable(t Table) string {
	return ""
}

func (driver DriverPostgres) RenameTable(t Table) string {
	return ""
}

func (driver DriverPostgres) DropTable(t Table) string {
	return ""
}

func (driver DriverPostgres) AddColumn(t Table, c TableColumn) string {
	return ""
}

func (driver DriverPostgres) DropColumn(t Table, c TableColumn) string {
	return ""
}

func (driver DriverPostgres) AlterColumn(t Table, c TableColumn) string {
	return ""
}

func (driver DriverPostgres) AddIndex(t Table, c TableIndex) string {
	return ""
}

func (driver DriverPostgres) DropIndex(t Table, c TableIndex) string {
	return ""
}

func (driver DriverPostgres) AddForeignKey(t Table, c TableForeignKey) string {
	return ""
}

func (driver DriverPostgres) DropForeignKey(t Table, c TableForeignKey) string {
	return ""
}

func (driver DriverPostgres) Insert(entity Entity) (string, map[string]any, error) {

	return "", nil, nil
}

func (driver DriverPostgres) Select(entity Entity) (Query, error) {
	return Query{}, nil
}

func (driver DriverPostgres) Update(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}

func (driver DriverPostgres) Save(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}

func (driver DriverPostgres) Delete(entity Entity) (string, map[string]any, error) {
	return "", nil, nil
}
