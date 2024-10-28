package database

type Driver interface {
	DSN() string
	Driver() string

	Insert(entity Entity) (string, map[string]any, error)
	Delete(entity Entity) (string, map[string]any, error)
	Save(entity Entity) (string, map[string]any, error)
	Select(entity Entity) (Query, error)
	Update(entity Entity) (string, map[string]any, error)

	TableFromEntity(e Entity) Table
	ParseCreateTable(s string) Table
	//
	CreateTable(t Table) string
	RenameTable() string
	DropTable(t Table) string
	//
	AddColumn(c TableColumn) string
	AlterColumn(c TableColumn) string
	DropColumn(c TableColumn) string
	//
	AddIndex(c TableIndex) string
	DropIndex(c TableIndex) string
	//
	AddForeignKey(c TableIndex) string
	DropForeignKey(c TableIndex) string
}
