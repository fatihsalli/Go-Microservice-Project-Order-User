package graphQL

import (
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/models"
	"context"
	"errors"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

var (
	orderType = graphql.NewObject(
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
			},
		},
	)

	rootQuery = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"orders": &graphql.Field{
					Type: graphql.NewList(orderType),
					Args: graphql.FieldConfigArgument{
						"status": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"userId": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						// Environment value
						env := os.Getenv("environment")

						// Get config
						config := configs.GetConfig(env)

						// Get collection
						mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.OrderCollectionName)

						statusArg, ok1 := params.Args["status"].(string)
						userArg, ok2 := params.Args["userId"].(string)

						if ok1 && ok2 {
							cursor, err := mongoOrderCollection.Find(context.Background(),
								bson.M{"status": statusArg, "userId": userArg})
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
			},
		},
	)

	Schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: rootQuery,
		},
	)
)
