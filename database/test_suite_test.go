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
	IsAdmin   bool       `db:"is_admin"`
	CreatedAt time.Time  `db:"created_at,readOnly,default=CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (u User1) EntityInformation() database.EntityInformation {
	return database.EntityInformation{
		TableName: "user",
	}
}

type User2 struct {
	ID        UserID    `db:"id,primaryKey,autoIncrement"`
	Name      string    `db:"name"`
	IsAdmin   bool      `db:"is_admin"`
	CreatedAt time.Time `db:"created_at,readOnly,default=CURRENT_TIMESTAMP"`
	// Missing DeletedAt
	Email EmailAddress `db:"email_address"` // New column
}

func (u User2) EntityInformation() database.EntityInformation {
	return database.EntityInformation{
		TableName: "user",
	}
}

func runDriverTestSuite(t *testing.T, dbal *database.DBAL) {
	t.Helper()

	ctx := context.Background()
	userRepository1 := database.NewRepository[UserID, User1](dbal)
	userRepository2 := database.NewRepository[UserID, User2](dbal)

	// Do the first migration
	if err := dbal.AutoMigrate(ctx, []database.Entity{
		User1{},
	}); err != nil {
		t.Fatal(err)
	}

	// Create a test user
	createdUser1ID, err := userRepository1.Insert(ctx, User1{
		Name:    "test user 1",
		IsAdmin: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Select the test user
	user1FromDB, err := userRepository1.GetByID(ctx, createdUser1ID)
	if err != nil {
		t.Fatal(err)
	}

	// Assert everything is working so far
	assert.Equal(t, UserID(1), createdUser1ID)
	assert.Equal(t, UserID(1), user1FromDB.ID)
	assert.Equal(t, true, user1FromDB.IsAdmin)

	// Do the second migration
	if err := dbal.AutoMigrate(ctx, []database.Entity{
		User2{},
	}); err != nil {
		t.Fatal(err)
	}

	// Update the existing one
	if err := userRepository2.Update(ctx, User2{
		ID:    createdUser1ID,
		Name:  "test user 2",
		Email: "foobar@example.com",
	}); err != nil {
		t.Fatal(err)
	}

	// Reselect the user
	userFromDB2, err := userRepository2.GetByID(ctx, createdUser1ID)
	if err != nil {
		t.Fatal(err)
	}

	// Assert all is working well after the second migration
	assert.Equal(t, UserID(1), userFromDB2.ID)
	assert.Equal(t, EmailAddress("foobar@example.com"), userFromDB2.Email)
}
