package mocks

import (
	"fmt"
	"sam-book-sample/models"

	"github.com/memememomo/dbmock"
)

func user(i uint64) dbmock.DBMapper {
	return &models.User{
		BaseModel: models.BaseModel{
			ID: i + 1,
		},
		Name:  fmt.Sprintf("Name_%d", i+1),
		Email: fmt.Sprintf("test_%d@example.com", i+1),
	}
}

func User() *dbmock.Generator {
	return dbmock.NewGenerator(user)
}
