package repository

import (
	"OrderUserProject/internal/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// TODO: Userid update edilmesi ayrı business
// TODO: Kafka da configleri collection gibi al
type OrderRepository struct {
	OrderCollection *mongo.Collection
}

func NewOrderRepository(mongoCollection *mongo.Collection) *OrderRepository {
	orderRepository := &OrderRepository{OrderCollection: mongoCollection}
	return orderRepository
}

// IOrderRepository to use for test or
type IOrderRepository interface {
	GetAll() ([]models.Order, error)
	GetOrderById(id string) (models.Order, error)
	Insert(order models.Order) (bool, error)
	Update(user models.Order) (bool, error)
	Delete(id string) (bool, error)
}

// GetAll Method => to list every order
func (b OrderRepository) GetAll() ([]models.Order, error) {
	var order models.Order
	var orders []models.Order

	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//We can think of "Cursor" like a request. We pull the data from the database with the "Next" command. (C# => IQueryable)
	result, err := b.OrderCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	for result.Next(ctx) {
		if err := result.Decode(&order); err != nil {
			return nil, err
		}
		// for appending book to books
		orders = append(orders, order)
	}

	return orders, nil

}

// GetOrderById Method => to find a single order with id
func (b OrderRepository) GetOrderById(id string) (models.Order, error) {
	var order models.Order

	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// to find book by id
	err := b.OrderCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)

	if err != nil {
		return order, err
	}

	return order, nil
}

// Insert method => to create new order
func (b OrderRepository) Insert(order models.Order) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// mongodb.driver
	result, err := b.OrderCollection.InsertOne(ctx, order)

	if result.InsertedID == nil || err != nil {
		return false, errors.New("failed to add")
	}

	return true, nil
}

// Update method => to change exist order
func (b OrderRepository) Update(order models.Order) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// => Update => update + insert = upsert => default value false
	// opt := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", order.ID}}

	// => if we use this CreatedDate and id value will be null, so we have to use "UpdateOne"
	//replacement := models.Book{Title: book.Title, Quantity: book.Quantity, Author: book.Author, UpdatedDate: book.UpdatedDate}

	// => to update for one parameter
	//update := bson.D{{"$set", bson.D{{"title", book.Title}}}}

	// => if we have to chance more than one parameter we have to write like this
	update := bson.D{{"$set", bson.D{
		{"userId", order.UserId},
		{"status", order.Status},
		{"address", order.Address},
		{"invoiceAddress", order.InvoiceAddress},
		{"product", order.Product},
		{"total", order.Total},
		{"updatedAt", order.UpdatedAt}}}}

	// mongodb.driver
	result, err := b.OrderCollection.UpdateOne(ctx, filter, update)

	if result.ModifiedCount <= 0 || err != nil {
		return false, err
	}

	return true, nil
}

// Delete Method => to delete a order from orders by id
func (b OrderRepository) Delete(id string) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// delete by id column
	result, err := b.OrderCollection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil || result.DeletedCount <= 0 {
		return false, err
	}

	return true, nil
}
