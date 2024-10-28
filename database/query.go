package database

import (
	"fmt"
	"strings"
)

func ParseQuery(query string) (Query, error) {
	return Query{}, nil
}

type Query struct {
	Select  []string
	From    string
	Joins   []string
	Where   string
	GroupBy string
	OrderBy string
	Limit   struct {
		Count  int
		Offset int
	}
	Parameters map[string]any
}

func (q Query) String() string {
	query := ""

	query += fmt.Sprintf("SELECT %s", strings.Join(q.Select, ", "))

	query += fmt.Sprintf(" FROM `%s`", q.From)

	query += strings.Join(q.Joins, " ")

	if q.Where != "" {
		query += fmt.Sprintf("WHERE %s", q.Where)
	}

	if q.GroupBy != "" {
		query += fmt.Sprintf("GROUP BY %s", q.GroupBy)
	}

	if q.OrderBy != "" {
		query += fmt.Sprintf("ORDER BY %s", q.OrderBy)
	}

	if q.Limit.Count != 0 && q.Limit.Offset != 0 {
		query += fmt.Sprintf("LIMIT %d,%d", q.Limit.Count, q.Limit)
	}

	return query
}
