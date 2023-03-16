package kafka

import (
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
	"log"
)

func ListenFromKafka(topic string) {

}

// SaveOrder for test to consume event and write on MongoDB
func SaveOrder(order models.Order) {
	//for test
	config := configs.GetConfig("test")
	mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection("Orders-event")
	OrderRepository := repository.NewOrderRepository(mongoOrderCollection)

	result, err := OrderRepository.Insert(order)
	if result == false || err != nil {
		log.Printf("Cannot create order event in MongoDB %v", order.ID)
	}
}
