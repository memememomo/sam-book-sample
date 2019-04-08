package models

import (
	"sam-book-sample/db"

	"github.com/memememomo/nomof"

	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

type UserEmailUniq struct {
	Email      string `dynamo:"PK"`
	EntityName string `dynamo:"SK"`
	Exists     bool   `dynamo:"Exists"`
	UserID     uint64 `dynamo:"UserID"`
}

func generateUserEmailUniqByUser(user *User) *UserEmailUniq {
	return &UserEmailUniq{
		Email:      user.Email,
		EntityName: getEntityNameFromStruct(*user),
		Exists:     true,
		UserID:     user.ID,
	}
}

func generateCreateQueryByUser(user *User) (*dynamo.Put, error) {
	table, err := db.Table()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	uniq := generateUserEmailUniqByUser(user)

	fb := nomof.NewBuilder()
	fb.AttributeNotExists("Exists")
	fb.Equal("UserID", user.ID)

	query := table.
		Put(&uniq).
		If(fb.JoinOr(), fb.Arg...)

	return query, nil
}

func generateDeleteQueryByUser(user *User) (*dynamo.Delete, error) {
	table, err := db.Table()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	uniq := generateUserEmailUniqByUser(user)

	query := table.
		Delete(db.PKName, uniq.Email).
		Range(db.SKName, uniq.EntityName)

	return query, nil
}
