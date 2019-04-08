package mocks

import (
	"sam-book-sample/db"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/memememomo/dbmock"
)

var E = func(i uint64, mapper dbmock.DBMapper) dbmock.DBMapper {
	return mapper
}

func SetupDB(t *testing.T) {
	t.Helper()
	err := db.SetupDBForTest()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func DumpDynamo(t *testing.T) {
	table, err := db.Table()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	var data []map[string]interface{}
	err = table.Scan().All(&data)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	pp.Print(data)
}
