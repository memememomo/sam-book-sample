package models

import (
	"sam-book-sample/db"

	"github.com/guregu/dynamo"

	"github.com/memememomo/nomof"
	"github.com/pkg/errors"
)

type Micropost struct {
	BaseModel
	Content string `dynamo:"Content"`
	UserID  uint64 `dynamo:"UserID"`
}

type MicropostDynamo struct {
	db.MainTable
	Micropost
}

// implements DynamoModelMapper

func (m *Micropost) EntityName() string {
	return getEntityNameFromStruct(*m)
}

func (m *Micropost) PK() string {
	return getPK(m)
}

func (m *Micropost) SK() string {
	return getSK(m)
}

func (m *Micropost) PutToDynamo() error {
	return putEntityToDynamo(m)
}

func (m *Micropost) SetID(id uint64) {
	m.ID = id
}

func (m *Micropost) GetID() uint64 {
	return m.ID
}

func (m *Micropost) SetVersion(v int) {
	m.Version = v
}

func (m *Micropost) GetVersion() int {
	return m.Version
}

// implements dbmock.ToDB

func (m *Micropost) ToDB() error {
	return m.Create()
}

func (m *Micropost) GenerateRecord() interface{} {
	return &MicropostDynamo{
		MainTable: generateMainTable(m),
		Micropost: *m,
	}
}

func (m *Micropost) CreateDynamoRecord() error {
	return createEntityToDynamo(m)
}

func (m *Micropost) UpdateDynamoRecord() error {
	return updateEntityToDynamo(m)
}

func (m *Micropost) DeleteDynamoRecord() error {
	return deleteEntity(m)
}

func (m *Micropost) Update() error {
	return m.PutToDynamo()
}

func (m *Micropost) Create() error {
	return m.PutToDynamo()
}

func GetMicropostByID(id uint64) (*Micropost, error) {
	var micropost MicropostDynamo
	ret, err := getEntityByID(id, &Micropost{}, &micropost)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if ret == nil {
		return nil, nil
	}
	return &micropost.Micropost, nil
}

func GetMicropostsByUserID(userID uint64) ([]*Micropost, error) {
	table, err := db.Table()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fb := nomof.NewBuilder()
	fb.Equal("UserID", userID)
	fb.BeginsWith(db.PKName, (&Micropost{}).EntityName())

	var micropostDynamo []MicropostDynamo
	err = table.
		Scan().
		Filter(fb.JoinAnd(), fb.Arg...).
		All(&micropostDynamo)

	if err != nil {
		if err == dynamo.ErrNotFound {
			return []*Micropost{}, nil
		}
		return nil, errors.WithStack(err)
	}

	var microposts = make([]*Micropost, len(micropostDynamo))
	for i := range micropostDynamo {
		microposts[i] = &micropostDynamo[i].Micropost
	}

	return microposts, nil
}

func DeleteMicropost(id uint64) error {
	micropost, err := GetMicropostByID(id)
	if err != nil {
		return errors.WithStack(err)
	}
	if micropost == nil {
		return nil
	}

	err = deleteEntity(micropost)

	return errors.WithStack(err)
}
