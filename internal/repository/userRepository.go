package repository

import (
	"OrderUserProject/internal/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepository struct {
	UserCollection *mongo.Collection
}

// TODO: singleton,scope ve transient olaylarına bakılacak

func NewUserRepository(mongoCollection *mongo.Collection) *UserRepository {
	userRepository := &UserRepository{UserCollection: mongoCollection}
	return userRepository
}

// IUserRepository to use for test or
type IUserRepository interface {
	GetAll() ([]models.User, error)
	GetUserById(id string) (models.User, error)
	Insert(user models.User) (bool, error)
	Update(user models.User) (bool, error)
	UpdateOrder(orderId string, userId string) (bool, error)
	Delete(id string) (bool, error)
	DeleteOrder(orderId string, userId string) (bool, error)
}

// GetAll Method => to list every user
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

// GetUserById Method => to find a single user with id
func (b UserRepository) GetUserById(id string) (models.User, error) {
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

// Insert method => to create new user
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

// Update method => to change exist user
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
		{"updateddate", user.UpdatedDate}}}}

	// mongodb.driver
	result, err := b.UserCollection.UpdateOne(ctx, filter, update)

	if result.ModifiedCount <= 0 || err != nil {
		return false, err
	}

	return true, nil
}

// UpdateOrder method => to change exist []string => Orders
func (b UserRepository) UpdateOrder(orderId string, userId string) (bool, error) {
	user, err := b.GetUserById(userId)

	if err != nil {
		return false, err
	}

	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// opt := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", user.ID}}

	// => if we have to chance more than one parameter we have to write like this
	update := bson.D{{"$set", bson.D{{"name", append(user.Orders, orderId)}}}}

	// mongodb.driver
	result, err := b.UserCollection.UpdateOne(ctx, filter, update)

	if result.ModifiedCount <= 0 || err != nil {
		return false, err
	}

	return true, nil
}

// Delete Method => to delete a user from users by id
func (b UserRepository) Delete(id string) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// delete by id column
	result, err := b.UserCollection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil || result.DeletedCount <= 0 {
		return false, err
	}

	return true, nil
}

// DeleteOrder method => to delete exist []string => Orders
func (b UserRepository) DeleteOrder(orderId string, userId string) (bool, error) {
	user, err := b.GetUserById(userId)

	if err != nil {
		return false, err
	}

	var newOrders []string

	for _, existOrderId := range user.Orders {
		if existOrderId != orderId {
			newOrders = append(newOrders, existOrderId)
		}
	}

	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// opt := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", user.ID}}

	// => if we have to chance more than one parameter we have to write like this
	update := bson.D{{"$set", bson.D{{"orders", newOrders}}}}

	// mongodb.driver
	result, err := b.UserCollection.UpdateOne(ctx, filter, update)

	if result.ModifiedCount <= 0 || err != nil {
		return false, err
	}

	return true, nil
}
