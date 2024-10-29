package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/lunagic/database-go/database"
	"gotest.tools/v3/assert"
)

type UserID uint64
type EmailAddress string

type User1 struct {
	ID        UserID     `db:"id,primaryKey,autoIncrement"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at,readOnly,default=CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (u User1) EntityInformation() database.EntityInformation {
	return database.EntityInformation{
		TableName: "user",
	}
}

type User2 struct {
	ID    UserID       `db:"id,primaryKey,autoIncrement"`
	Name  string       `db:"name"`
	Email EmailAddress `db:"email_address"` // New column
	// Missing UpdatedAt
}

func (u User2) EntityInformation() database.EntityInformation {
	return database.EntityInformation{
		TableName: "user",
	}
}

func runDriverTestSuite(t *testing.T, dbal *database.DBAL) {
	ctx := context.Background()

	if err := dbal.AutoMigrate(ctx, []database.Entity{
		User1{},
	}); err != nil {
		t.Fatal(err)
	}

	userRepository1 := database.NewRepository[UserID, User1](dbal)
	userRepository2 := database.NewRepository[UserID, User2](dbal)

	createdUserID, err := userRepository1.Insert(ctx, User1{
		Name: "test user 1",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, UserID(1), createdUserID)

	userFromDB, err := userRepository1.SelectSingle(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, UserID(1), userFromDB.ID)

	if err := dbal.AutoMigrate(ctx, []database.Entity{
		User2{},
	}); err != nil {
		t.Fatal(err)
	}

	userRepository2.Update(ctx, User2{
		ID:    createdUserID,
		Name:  "test user 2",
		Email: "foobar@example.com",
	})

	userFromDB2, err := userRepository2.SelectSingle(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, UserID(1), userFromDB2.ID)
	assert.Equal(t, EmailAddress("foobar@example.com"), userFromDB2.Email)

}
