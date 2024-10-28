package database

import "reflect"

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
	Type          reflect.Type
	AutoIncrement bool
	PrimaryKey    bool
	Nullable      string
	Index         []TableIndex
	ForeignKey    *TableForeignKey
}

func EntityToTable(entity Entity) (Table, error) {
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
			Type:          fieldDefinition.Type,
			PrimaryKey:    tag.PrimaryKey,
		}

		table.Columns = append(table.Columns, column)
	})

	return table, nil
}

func (dbal *DBAL) diffFunc(a Table, b Table) []string {
	return []string{}
}
