package database

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func NewDBAL(driver Driver, configFuncs ...ConfigFunc) (*DBAL, error) {
	connection, err := sql.Open(driver.Driver(), driver.DSN())
	if err != nil {
		return nil, err
	}

	_ = mysql.SetLogger(log.New(io.Discard, "", log.LstdFlags))

	dbal := &DBAL{
		connection: connection,
		driver:     driver,
	}

	for _, configFunc := range configFuncs {
		if err := configFunc(dbal); err != nil {
			return nil, err
		}
	}

	return dbal, nil
}

type ConfigPreRunFunc func(ctx context.Context, statement string, args []any) error

type ConfigPostRunFunc func(ctx context.Context) error

type DBAL struct {
	connection   *sql.DB
	driver       Driver
	preRunFuncs  []ConfigPreRunFunc
	postRunFuncs []ConfigPostRunFunc
}

type tableCreateTargetStruct struct {
	Table       string `db:"Table"`
	CreateTable string `db:"Create Table"`
}

func (dbal *DBAL) AutoMigrate(ctx context.Context, entities []Entity) error {
	for _, entity := range entities {
		targetTableDefinition := dbal.driver.TableFromEntity(entity)

		tableCreateTarget := []tableCreateTargetStruct{}
		if err := dbal.RawSelect(ctx, dbal.driver.ShowCreateTable(entity.EntityInformation().TableName), nil, &tableCreateTarget); err != nil {
			if !strings.Contains(err.Error(), "42S02") {
				return err
			}
		}

		if len(tableCreateTarget) < 1 {
			if _, err := dbal.RawExecute(ctx, dbal.driver.CreateTable(targetTableDefinition), nil); err != nil {
				return err
			}

			continue
		}

		currentTableDefinition := dbal.driver.ParseCreateTable(tableCreateTarget[0].CreateTable)
		currentTableDefinition.Name = targetTableDefinition.Name // TODO: remove this

		updateStatements := dbal.diffFunc(targetTableDefinition, currentTableDefinition)

		for _, updateStatement := range updateStatements {
			if _, err := dbal.RawExecute(ctx, updateStatement, nil); err != nil {
				return err
			}
		}

		return nil
	}

	// destTable := dbal.driver.getCreateSyntax(entities[0].EntityInformation().TableName)
	// srcTable := TableFromEntity(entities[0])
	// diff := Diff(destTable, srcTable)
	// dbal.RawExecute(ctx, diff[0])
	return nil
}

func (dbal *DBAL) RawSelect(
	ctx context.Context,
	query string,
	parameters map[string]any,
	targetPointer any,
) error {
	preparedQuery, preparedArgs, err := prepare(query, parameters)
	if err != nil {
		return err
	}

	for _, preRunFunc := range dbal.preRunFuncs {
		if err := preRunFunc(ctx, preparedQuery, preparedArgs); err != nil {
			return err
		}
	}

	rows, err := dbal.connection.Query(preparedQuery, preparedArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	target := reflect.ValueOf(targetPointer).Elem()
	targetType := reflect.TypeOf(target.Interface()).Elem()

	fieldIndexesToUse := []int{}

	rowMap := map[string]int{}
	testRow := reflect.New(targetType).Elem()
	for i := 0; i < testRow.NumField(); i++ {
		fieldDefinition := testRow.Type().Field(i)
		if !fieldDefinition.IsExported() {
			continue
		}

		tag := parseTag(fieldDefinition.Tag)
		if tag.Column == "" {
			continue
		}

		rowMap[tag.Column] = i
	}

	for _, column := range columns {
		fieldIndex, found := rowMap[column]
		if !found {
			return fmt.Errorf("column %s not found in target", column)
		}

		fieldIndexesToUse = append(fieldIndexesToUse, fieldIndex)
	}

	for rows.Next() {
		row := reflect.New(targetType).Elem()

		scanFields := []any{}
		for _, fieldIndexToUse := range fieldIndexesToUse {
			scanFields = append(scanFields, row.Field(fieldIndexToUse).Addr().Interface())
		}

		if err := rows.Scan(scanFields...); err != nil {
			return err
		}

		target.Set(reflect.Append(target, row))
	}

	for _, postRunFunc := range dbal.postRunFuncs {
		if err := postRunFunc(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (dbal *DBAL) RawExecute(
	ctx context.Context,
	query string,
	parameters map[string]any,
) (sql.Result, error) {
	preparedQuery, preparedArgs, err := prepare(query, parameters)
	if err != nil {
		return nil, err
	}

	for _, preRunFunc := range dbal.preRunFuncs {
		if err := preRunFunc(ctx, preparedQuery, preparedArgs); err != nil {
			return nil, err
		}
	}

	result, err := dbal.connection.Exec(preparedQuery, preparedArgs...)
	if err != nil {
		return nil, err
	}

	for _, postRunFunc := range dbal.postRunFuncs {
		if err := postRunFunc(ctx); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (dbal *DBAL) Ping() error {
	return dbal.connection.Ping()
}
