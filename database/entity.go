package database

type EntityInformation struct {
	TableName    string
	Indexes      []TableIndex
	ForeignKey   []TableForeignKey
	OldTableName string
}

type Entity interface {
	EntityInformation() EntityInformation
}

type TableIndex struct{}

type TableForeignKey struct{}

func DiffFunc(a Table, b Table) []string {
	return []string{}
}

type Table struct {
	Columns TableColumn
}

type TableColumn struct {
	Name     string
	Type     string // int(20) unsigned -  uint64
	Position int
	Comment  string
	Nullable bool
}
