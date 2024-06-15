package database

import (
	"context"
)

func NewRepository[ID ~uint64, T Entity](dbal *DBAL) Repository[ID, T] {
	query, err := GenerateSelect(*new(T))
	if err != nil {
		panic(err)
	}

	return Repository[ID, T]{
		Selector: NewSelector[T](dbal, query),
		dbal:     dbal,
	}
}

type Repository[ID ~uint64, T Entity] struct {
	Selector[T]
	dbal *DBAL
}

func (q *Repository[ID, T]) Insert(ctx context.Context, entity T) (ID, error) {
	statement, parameters, err := GenerateInsert(entity)
	if err != nil {
		return 0, err
	}

	result, err := q.connection.RawExecute(ctx, statement, parameters)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return ID(lastInsertID), nil
}

func (q *Repository[ID, T]) Update(ctx context.Context, entity T) error {
	statement, parameters, err := GenerateUpdate(entity)
	if err != nil {
		return err
	}

	if _, err := q.connection.RawExecute(ctx, statement, parameters); err != nil {
		return err
	}

	return nil
}

func (q *Repository[ID, T]) Save(ctx context.Context, entity T) error {
	statement, parameters, err := GenerateSave(entity)
	if err != nil {
		return err
	}

	if _, err := q.connection.RawExecute(ctx, statement, parameters); err != nil {
		return err
	}

	return nil
}
