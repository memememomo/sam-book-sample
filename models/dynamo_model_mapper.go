package models

import (
	"fmt"
	"reflect"
	"sam-book-sample/db"

	"github.com/guregu/dynamo"
	"github.com/memememomo/nomof"
	"github.com/pkg/errors"
)

type DynamoModelMapper interface {
	EntityName() string
	PK() string
	SK() string
	PutToDynamo() error
	SetID(id uint64)
	GetID() uint64
	SetVersion(v int)
	GetVersion() int
	GenerateRecord() interface{}
}

type DynamoEntityMapper interface {
	DynamoModelMapper
	CreateDynamoRecord() error
	UpdateDynamoRecord() error
	DeleteDynamoRecord() error
}

func getEntityNameFromStruct(s interface{}) string {
	r := reflect.TypeOf(s)
	return r.Name()
}

func getPK(s DynamoModelMapper) string {
	return fmt.Sprintf("%s-%011d", s.EntityName(), s.GetID())
}

func getSK(s DynamoModelMapper) string {
	return fmt.Sprintf("%011d", s.GetID())
}

func generateMainTable(mapper DynamoModelMapper) db.MainTable {
	return db.MainTable{
		PK: mapper.PK(),
		SK: mapper.SK(),
	}
}

func generateCreateQuery(mapper DynamoEntityMapper) (*dynamo.Put, error) {
	table, err := db.Table()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	id, err := db.GenerateID(mapper.EntityName())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	mapper.SetID(id)
	mapper.SetVersion(1)

	fb := nomof.NewBuilder()
	fb.AttributeNotExists(db.PKName)

	query := table.
		Put(mapper.GenerateRecord()).
		If(fb.JoinAnd(), fb.Arg...)

	return query, nil
}

func generateUpdateQuery(mapper DynamoEntityMapper) (*dynamo.Put, error) {
	table, err := db.Table()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	oldVersion := mapper.GetVersion()

	mapper.SetVersion(oldVersion + 1)

	fb := nomof.NewBuilder()
	fb.Equal("Version", oldVersion)

	query := table.
		Put(mapper.GenerateRecord()).
		If(fb.JoinAnd(), fb.Arg...)

	return query, nil
}

func generateDeleteQuery(mapper DynamoEntityMapper) (*dynamo.Delete, error) {
	table, err := db.Table()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	query := table.
		Delete(db.PKName, mapper.PK()).
		Range(db.SKName, mapper.SK())

	return query, nil
}

func createEntityToDynamo(mapper DynamoEntityMapper) error {
	query, err := generateCreateQuery(mapper)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func updateEntityToDynamo(mapper DynamoEntityMapper) error {
	query, err := generateUpdateQuery(mapper)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func deleteEntity(mapper DynamoEntityMapper) error {
	query, err := generateDeleteQuery(mapper)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func isNewEntity(mapper DynamoEntityMapper) bool {
	return mapper.GetVersion() == 0
}

func putEntityToDynamo(mapper DynamoEntityMapper) error {
	if isNewEntity(mapper) {
		return mapper.CreateDynamoRecord()
	}
	return mapper.UpdateDynamoRecord()
}

func getEntityByID(id uint64, mapper DynamoEntityMapper, ret interface{}) (interface{}, error) {
	table, err := db.Table()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	mapper.SetID(id)
	err = table.
		Get(db.PKName, mapper.PK()).
		Range(db.SKName, dynamo.Equal, mapper.SK()).
		One(ret)

	if err != nil {
		if err == dynamo.ErrNotFound {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	return ret, nil
}
