package database

type Driver interface {
	// Connection Methods
	Driver() string
	DSN() string
	// Migration Methods
	ShowCreateTable(tableName string) string
	ParseCreateTable(s string) Table
	TableFromEntity(e Entity) Table
	CreateTable(t Table) string
	RenameTable(t Table) string
	DropTable(t Table) string
	AddColumn(t Table, c TableColumn) string
	AlterColumn(t Table, c TableColumn) string
	DropColumn(t Table, c TableColumn) string
	AddIndex(t Table, c TableIndex) string
	DropIndex(t Table, c TableIndex) string
	AddForeignKey(t Table, c TableForeignKey) string
	DropForeignKey(t Table, c TableForeignKey) string
	// CRUD Methods
	Insert(e Entity) (string, map[string]any, error)
	Update(e Entity) (string, map[string]any, error)
	Save(e Entity) (string, map[string]any, error)
	Select(e Entity) (Query, error)
	Delete(e Entity) (string, map[string]any, error)
}
