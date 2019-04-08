package mocks

import (
	"fmt"
	"sam-book-sample/models"

	"github.com/memememomo/dbmock"
)

func micropost(i uint64) dbmock.DBMapper {
	return &models.Micropost{
		BaseModel: models.BaseModel{
			ID: i + 1,
		},
		UserID:  1,
		Content: fmt.Sprintf("Content_%d", i+1),
	}
}

func Micropost() *dbmock.Generator {
	return dbmock.NewGenerator(micropost)
}
