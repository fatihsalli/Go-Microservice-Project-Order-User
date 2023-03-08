package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection  *mongo.Collection
	constructor func() interface{}
}

func NewRepository(coll *mongo.Collection, cons func() interface{}) *Repository {
	return &Repository{
		collection:  coll,
		constructor: cons,
	}
}

func (r *Repository) GetAll(ctx context.Context) ([]interface{}, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var result []interface{}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		entry := r.constructor() // call to constructor
		if err = cur.Decode(entry); err != nil {
			return nil, err
		}

		result = append(result, entry)
	}

	return result, nil
}

func (r *Repository) Create(ctx context.Context, obj interface{}) error {
	_, err := r.collection.InsertOne(ctx, obj)
	if err != nil {
		return err
	}

	return nil
}
