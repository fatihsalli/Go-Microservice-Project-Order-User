package repository

import (
	"OrderUserProject/internal/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepository struct {
	UserCollection *mongo.Collection
}

var singleInstanceRepo *UserRepository

// TODO: singleton,scope ve transient olaylarına bakılacak
func GetSingleInstancesRepository(mongoCollection *mongo.Collection) *UserRepository {
	if singleInstanceRepo == nil {
		fmt.Println("Creating single repository instance now.")
		singleInstanceRepo = &UserRepository{UserCollection: mongoCollection}
	} else {
		fmt.Println("Single repository instance already created.")
	}

	return singleInstanceRepo
}

// IUserRepository to use for test or
type IUserRepository interface {
	Insert(user models.User) (bool, error)
	GetAll() ([]models.User, error)
	GetBookById(id string) (models.User, error)
	Update(user models.User) (bool, error)
	Delete(id string) (bool, error)
}

// Insert method => to create new book
func (b UserRepository) Insert(user models.User) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// mongodb.driver
	result, err := b.UserCollection.InsertOne(ctx, user)

	if result.InsertedID == nil || err != nil {
		return false, errors.New("failed to add")
	}

	return true, nil
}

// Update method => to change exist book
func (b UserRepository) Update(user models.User) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// => Update => update + insert = upsert => default value false
	// opt := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", user.ID}}

	// => if we use this CreatedDate and id value will be null, so we have to use "UpdateOne"
	//replacement := models.Book{Title: book.Title, Quantity: book.Quantity, Author: book.Author, UpdatedDate: book.UpdatedDate}

	// => to update for one parameter
	//update := bson.D{{"$set", bson.D{{"title", book.Title}}}}

	// => if we have to chance more than one parameter we have to write like this
	update := bson.D{{"$set", bson.D{{"name", user.Name},
		{"email", user.Email}, {"password", user.Password}, {"address", user.Address},
		{"orders", user.Orders}, {"updated_date", user.UpdatedDate}}}}

	// mongodb.driver
	result, err := b.UserCollection.UpdateOne(ctx, filter, update)

	if result.ModifiedCount <= 0 || err != nil {
		return false, err
	}

	return true, nil
}

// GetAll Method => to list every books
func (b UserRepository) GetAll() ([]models.User, error) {
	var user models.User
	var users []models.User

	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//We can think of "Cursor" like a request. We pull the data from the database with the "Next" command. (C# => IQueryable)
	result, err := b.UserCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	for result.Next(ctx) {
		if err := result.Decode(&user); err != nil {
			return nil, err
		}
		// for appending book to books
		users = append(users, user)
	}

	return users, nil

}

// GetBookById Method => to find a single book with id
func (b UserRepository) GetBookById(id string) (models.User, error) {
	var user models.User

	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// to find book by id
	err := b.UserCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	if err != nil {
		return user, err
	}

	return user, nil
}

// Delete Method => to delete a book from books by id
func (b UserRepository) Delete(id string) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// delete by id column
	result, err := b.UserCollection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil || result.DeletedCount <= 0 {
		return false, err
	}

	return true, nil
}
