package database

type EntityInformation struct {
	TableName  string
	Indexes    []TableIndex
	ForeignKey []TableForeignKey
	// OldTableName string // TODO: use for possibly support table renames
}

type Entity interface {
	EntityInformation() EntityInformation
}

type TableIndex struct{}

type TableForeignKey struct{}

type Table struct {
	Name    string
	Columns []TableColumn
}

type TableColumn struct {
	Name          string
	Type          string
	AutoIncrement bool
	PrimaryKey    bool
	HasDefault    bool
	Default       string
	Nullable      bool
	Index         []TableIndex
	ForeignKey    *TableForeignKey
}

func (dbal *DBAL) diffFunc(targetTable Table, currentTable Table) []string {

	currentColumns := map[string]TableColumn{}
	for _, column := range currentTable.Columns {
		currentColumns[column.Name] = column
	}

	targetColumns := map[string]TableColumn{}
	for _, column := range targetTable.Columns {
		targetColumns[column.Name] = column
	}

	columnsToAdd := []string{}
	columnsToDrop := []string{}
	columnsToAlter := []string{}

	for columnName, column := range currentColumns {
		if _, found := targetColumns[columnName]; !found {
			columnsToDrop = append(
				columnsToDrop,
				dbal.driver.DropColumn(targetTable, column),
			)
		}
	}

	for columnName, column := range targetColumns {
		if _, found := currentColumns[columnName]; !found {
			columnsToAdd = append(
				columnsToAdd,
				dbal.driver.AddColumn(targetTable, column),
			)
		}
	}

	return append(columnsToDrop, append(columnsToAdd, columnsToAlter...)...)
}
