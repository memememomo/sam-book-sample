package models

import "time"

type CreatedUpdated struct {
	CreatedAt time.Time `dynamo:"CreatedAt"`
	UpdatedAt time.Time `dynamo:"UpdatedAt"`
}

type BaseModel struct {
	ID      uint64 `dynamo:"ID"`
	Version int    `dynamo:"Version"`
	CreatedUpdated
}
