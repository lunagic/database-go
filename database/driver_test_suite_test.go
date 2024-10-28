package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/lunagic/database-go/database"
	"github.com/stretchr/testify/assert"
)

type UserID uint64

type User1 struct {
	ID        UserID    `db:"id,primaryKey"`
	Name      string    `db:"name"`
	UpdatedAt time.Time `db:"updated_at,readOnly"`
}

func (u User1) EntityInformation() database.EntityInformation {
	return database.EntityInformation{
		TableName: "user",
	}
}

type User1Repository struct {
	database.Repository[UserID, User1]
}

func runDriverTestSuite(t *testing.T, dbal *database.DBAL) {
	ctx := context.Background()

	if err := dbal.AutoMigrate(ctx, []database.Entity{
		User1{},
	}); err != nil {
		t.Fatal(err)
	}

	userRepository := User1Repository{
		Repository: database.NewRepository[UserID, User1](dbal),
	}

	createdUserID, err := userRepository.Insert(ctx, User1{})
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, UserID(1), createdUserID) {
		return
	}

	userFromDB, err := userRepository.SelectSingle(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, UserID(1), userFromDB.ID) {
		return
	}
}
