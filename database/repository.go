package database

import (
	"context"
)

func NewRepository[ID ~uint64, T Entity](dbal *DBAL) Repository[ID, T] {
	query, err := dbal.driver.Select(*new(T))
	if err != nil {
		panic(err)
	}

	return Repository[ID, T]{
		Selector: NewSelector[T](dbal, query),
	}
}

type Repository[ID ~uint64, T Entity] struct {
	Selector[T]
}

func (q *Repository[ID, T]) GetByID(ctx context.Context, id ID) (T, error) {
	return q.SelectSingle(ctx, WithAdditionalWhere("`id` = :id", map[string]any{
		":id": id,
	}))
}

func (q *Repository[ID, T]) Insert(ctx context.Context, entity T) (ID, error) {
	statement, parameters, err := q.Selector.dbal.driver.Insert(entity)
	if err != nil {
		return 0, err
	}

	result, err := q.dbal.RawExecute(ctx, statement, parameters)
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
	statement, parameters, err := q.dbal.driver.Update(entity)
	if err != nil {
		return err
	}

	if _, err := q.dbal.RawExecute(ctx, statement, parameters); err != nil {
		return err
	}

	return nil
}

func (q *Repository[ID, T]) Save(ctx context.Context, entity T) error {
	statement, parameters, err := q.Selector.dbal.driver.Save(entity)
	if err != nil {
		return err
	}

	if _, err := q.dbal.RawExecute(ctx, statement, parameters); err != nil {
		return err
	}

	return nil
}
