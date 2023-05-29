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
	addressType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Address",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"address": &graphql.Field{
				Type: graphql.String,
			},
			"city": &graphql.Field{
				Type: graphql.String,
			},
			"district": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"default": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "AddressDefault",
					Fields: graphql.Fields{
						"isDefaultInvoiceAddress": &graphql.Field{
							Type: graphql.Boolean,
						},
						"isDefaultRegularAddress": &graphql.Field{
							Type: graphql.Boolean,
						},
					},
				}),
			},
		},
	})

	orderType = graphql.NewObject(graphql.ObjectConfig{
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
			"address": &graphql.Field{
				Type: addressType,
			},
			"invoiceAddress": &graphql.Field{
				Type: addressType,
			},
			"product": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
					Name: "Product",
					Fields: graphql.Fields{
						"name": &graphql.Field{
							Type: graphql.String,
						},
						"quantity": &graphql.Field{
							Type: graphql.Int,
						},
						"price": &graphql.Field{
							Type: graphql.Float,
						},
					},
				})),
			},
			"total": &graphql.Field{
				Type: graphql.Float,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	})

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
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						// Environment value
						env := os.Getenv("environment")

						// Get config
						config := configs.GetConfig(env)

						// Get collection
						mongoOrderCollection := configs.ConnectDB(config.Database.Connection).
							Database(config.Database.DatabaseName).
							Collection(config.Database.OrderCollectionName)

						statusArg, ok := params.Args["status"].(string)

						if ok {
							cursor, err := mongoOrderCollection.Find(context.Background(),
								bson.M{"status": statusArg})
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
