package database

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Hostname string
	Port     int
	Username string
	Password string
	Name     string
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?parseTime=true",
		c.Username,
		c.Password,
		c.Hostname,
		c.Port,
		c.Name,
	)
}

func NewDBAL(config Config, logger *log.Logger) (*DBAL, error) {
	connection, err := sql.Open("mysql", config.DSN())
	if err != nil {
		return nil, err
	}

	mysql.SetLogger(log.New(io.Discard, "", log.LstdFlags))

	return &DBAL{
		connection: connection,
		logger:     logger,
	}, nil
}

type DBAL struct {
	connection *sql.DB
	logger     *log.Logger
}

func (dbal *DBAL) RawSelect(
	ctx context.Context,
	query string,
	parameters map[string]any,
	targetPointer any,
) error {
	preparedQuery, preparedArgs, err := Prepare(query, parameters)
	if err != nil {
		return err
	}

	dbal.logf("DBAL Select: %s", query)

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

	return nil
}

func (dbal *DBAL) RawExecute(
	ctx context.Context,
	query string,
	parameters map[string]any,
) (sql.Result, error) {
	preparedQuery, preparedArgs, err := Prepare(query, parameters)
	if err != nil {
		return nil, err
	}

	dbal.logf("DBAL Execute: %s", query)

	return dbal.connection.Exec(preparedQuery, preparedArgs...)
}

func (dbal *DBAL) Ping() error {
	return dbal.connection.Ping()
}

func (dbal *DBAL) logf(
	s string,
	args ...any,
) {
	if dbal.logger == nil {
		return
	}

	dbal.logger.Printf(s, args...)

}
