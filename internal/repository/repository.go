package repository

import "go.mongodb.org/mongo-driver/mongo"

type repository[T any] struct {
	db *mongo.Database
}

func NewRepository[T any](db *mongo.Database) *repository[T] {
	return &repository[T]{
		db: db,
	}
}
