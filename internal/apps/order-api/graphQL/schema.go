package graphQL

import (
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

var orderType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Order",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"userId": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"total": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) single product by id
			   http://localhost:30011/order?query={order(id:"e9caaa02-5c6a-4d2f-b795-11680de70401"){userId,status,total}}
			*/
			"order": &graphql.Field{
				Type:        orderType,
				Description: "Get order by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Environment value
					env := os.Getenv("environment")

					// Get config
					config := configs.GetConfig(env)

					// Get collection
					mongoOrderCollection := configs.ConnectDB(config.Database.Connection).
						Database(config.Database.DatabaseName).
						Collection(config.Database.OrderCollectionName)

					id, ok := p.Args["id"].(string)

					if ok {
						cursor, err := mongoOrderCollection.Find(context.Background(),
							bson.M{"_id": id})
						if err != nil {
							return nil, err
						}
						defer cursor.Close(context.Background())

						var orders []models.Order
						for cursor.Next(context.Background()) {
							var order models.Order
							err := cursor.Decode(&order)
							if err != nil {
								return nil, err
							}
							orders = append(orders, order)
						}

						return orders, nil
					} else {
						return nil, errors.New("bad request for graphQL")
					}
				},
			},
			/* Get (read) product list
			   http://localhost:30011/order?query={list{id,userId,status,total}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(orderType),
				Description: "Get order list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Environment value
					env := os.Getenv("environment")

					// Get config
					config := configs.GetConfig(env)

					// Get collection
					mongoOrderCollection := configs.ConnectDB(config.Database.Connection).
						Database(config.Database.DatabaseName).
						Collection(config.Database.OrderCollectionName)

					cursor, err := mongoOrderCollection.Find(context.Background(),
						bson.M{})
					if err != nil {
						return nil, err
					}
					defer cursor.Close(context.Background())

					var orders []models.Order
					for cursor.Next(context.Background()) {
						var order models.Order
						err := cursor.Decode(&order)
						if err != nil {
							return nil, err
						}
						orders = append(orders, order)
					}

					return orders, nil
				},
			},
		},
	})

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func ExecuteQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}
