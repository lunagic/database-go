package database_test

import (
	"testing"
	"time"

	"github.com/lunagic/database-go/database"
	"github.com/stretchr/testify/assert"
)

type UserID uint64

type User struct {
	ID        UserID    `db:"id,primaryKey"`
	Name      string    `db:"name"`
	UpdatedAt time.Time `db:"updated_at,readOnly"`
}

func (u User) EntityInformation() database.EntityInformation {
	return database.EntityInformation{
		TableName: "user",
	}
}

type UserRepo struct {
	database.Repository[UserID, User]
}

func TestFoobar(t *testing.T) {
	ctx, dbal := getDockerDBAL(t)

	if _, err := dbal.RawExecute(ctx, `
		CREATE TABLE user (
			id int(11) NOT NULL AUTO_INCREMENT,
			name varchar(255) DEFAULT NULL,
			updated_at timestamp NULL DEFAULT current_timestamp(),
			PRIMARY KEY (id)
		);
	`, nil); err != nil {
		t.Fatal(err)
	}

	r := UserRepo{
		Repository: database.NewRepository[UserID, User](dbal),
	}

	createdUserID, err := r.Insert(ctx, User{})
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, UserID(1), createdUserID) {
		return
	}

	userFromDB, err := r.SelectSingle(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, UserID(1), userFromDB.ID) {
		return
	}
}
