package database

import (
	"context"
	"errors"
)

var ErrNoRows = errors.New("no rows")

func NewSelector[T any](connection *DBAL, query Query) Selector[T] {
	return Selector[T]{
		connection: connection,
		query:      query,
	}
}

type Selector[T any] struct {
	connection *DBAL
	query      Query
}

type QueryModifier func(query Query) Query

func WithLimitOverride(size int, offset int) QueryModifier {
	return func(query Query) Query {
		query.Limit.Count = size
		query.Limit.Offset = offset

		return query
	}
}

func (q *Selector[T]) SelectMultiple(ctx context.Context, mods ...QueryModifier) ([]T, error) {
	target := []T{}

	query := q.query
	for _, mod := range mods {
		query = mod(query)
	}

	if err := q.connection.RawSelect(ctx, q.query.String(), query.Parameters, &target); err != nil {
		return nil, err
	}

	return target, nil
}

func (q *Selector[T]) SelectSingle(ctx context.Context, mods ...QueryModifier) (T, error) {
	mods = append(mods, WithLimitOverride(1, 0))

	rows, err := q.SelectMultiple(ctx, mods...)
	if err != nil {
		return *new(T), err
	}

	if len(rows) < 1 {
		return *new(T), ErrNoRows
	}

	return rows[0], nil
}
