package models

import (
	"sam-book-sample/db"

	"github.com/guregu/dynamo"

	"github.com/memememomo/nomof"
	"github.com/pkg/errors"
)

type User struct {
	BaseModel
	Name  string `dynamo:"Name"`
	Email string `dynamo:"Email"`
}

type UserDynamo struct {
	db.MainTable
	User
}

// implements DynamoModelMapper

func (u *User) EntityName() string {
	return getEntityNameFromStruct(*u)
}

func (u *User) PK() string {
	return getPK(u)
}

func (u *User) SK() string {
	return getSK(u)
}

func (u *User) PutToDynamo() error {
	return putEntityToDynamo(u)
}

func (u *User) CreateDynamoRecord() error {
	conn, err := db.ConnectDB()
	if err != nil {
		return errors.WithStack(err)
	}

	tx := conn.WriteTx()

	r, err := generateCreateQuery(u)
	if err != nil {
		return errors.WithStack(err)
	}

	uniq, err := generateCreateQueryByUser(u)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Put(r).Put(uniq).Run()

	return errors.WithStack(err)
}

func (u *User) UpdateDynamoRecord() error {
	conn, err := db.ConnectDB()
	if err != nil {
		return errors.WithStack(err)
	}

	tx := conn.WriteTx()

	r, err := generateUpdateQuery(u)
	if err != nil {
		return errors.WithStack(err)
	}

	uniq, err := generateCreateQueryByUser(u)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Put(r).Put(uniq).Run()

	return errors.WithStack(err)
}

func (u *User) DeleteDynamoRecord() error {
	conn, err := db.ConnectDB()
	if err != nil {
		return errors.WithStack(err)
	}

	tx := conn.WriteTx()

	r, err := generateDeleteQuery(u)
	if err != nil {
		return errors.WithStack(err)
	}

	uniq, err := generateDeleteQueryByUser(u)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Delete(r).Delete(uniq).Run()

	return errors.WithStack(err)
}

func (u *User) SetID(id uint64) {
	u.ID = id
}

func (u *User) GetID() uint64 {
	return u.ID
}

func (u *User) SetVersion(v int) {
	u.Version = v
}

func (u *User) GetVersion() int {
	return u.Version
}

// implement dbmock.DBMapper

func (u *User) ToDB() error {
	return u.Create()
}

func (u *User) GenerateRecord() interface{} {
	return &UserDynamo{
		MainTable: generateMainTable(u),
		User:      *u,
	}
}

func (u *User) Update() error {
	return u.PutToDynamo()
}

func (u *User) Create() error {
	return u.PutToDynamo()
}

func GetUserByEmail(email string) (*User, error) {
	table, err := db.Table()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fb := nomof.NewBuilder()
	fb.Equal("Email", email)
	fb.BeginsWith(db.PKName, (&User{}).EntityName())

	var usersDynamo []UserDynamo
	err = table.
		Scan().
		Filter(fb.JoinAnd(), fb.Arg...).
		All(&usersDynamo)

	if err != nil {
		if err == dynamo.ErrNotFound {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	if len(usersDynamo) == 0 {
		return nil, nil
	}

	return &usersDynamo[0].User, nil
}

func GetUserByID(id uint64) (*User, error) {
	var user UserDynamo
	ret, err := getEntityByID(id, &User{}, &user)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if ret == nil {
		return nil, nil
	}
	return &user.User, nil
}

func GetUsers() ([]*User, error) {
	table, err := db.Table()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fb := nomof.NewBuilder()
	fb.BeginsWith(db.PKName, (&User{}).EntityName())

	var userDynamo []UserDynamo
	err = table.
		Scan().
		Filter(fb.JoinAnd(), fb.Arg...).
		All(&userDynamo)

	if err != nil {
		if err == dynamo.ErrNotFound {
			return []*User{}, nil
		}
		return nil, errors.WithStack(err)
	}

	var users = make([]*User, len(userDynamo))
	for i := 0; i < len(userDynamo); i++ {
		users[i] = &userDynamo[i].User
	}

	return users, nil
}

func DeleteUser(id uint64) error {
	user, err := GetUserByID(id)
	if err != nil {
		return errors.WithStack(err)
	}
	if user == nil {
		return nil
	}

	err = deleteEntity(user)

	return errors.WithStack(err)
}
