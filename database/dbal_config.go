package database

import (
	"context"
	"encoding/json"
	"log"
)

type ConfigFunc func(dbal *DBAL) error

func WithPreRunFunc(preRunFunc ConfigPreRunFunc) ConfigFunc {
	return func(dbal *DBAL) error {
		dbal.preRunFuncs = append(dbal.preRunFuncs, preRunFunc)
		return nil
	}
}

func WithPostRunFunc(postRunFunc ConfigPostRunFunc) ConfigFunc {
	return func(dbal *DBAL) error {
		dbal.postRunFuncs = append(dbal.postRunFuncs, postRunFunc)
		return nil
	}
}

func WithLogger(logger *log.Logger) ConfigFunc {
	return func(dbal *DBAL) error {
		dbal.preRunFuncs = append(dbal.preRunFuncs, func(ctx context.Context, statement string, args []any) error {
			argJSON, err := json.Marshal(args)
			if err != nil {
				return err
			}

			logger.Printf("DBAL Run: %s %s", statement, string(argJSON))

			return nil
		})
		dbal.postRunFuncs = append(dbal.postRunFuncs, func(ctx context.Context) error {
			// TODO: log how long it took to complete?
			return nil
		})
		return nil
	}
}
