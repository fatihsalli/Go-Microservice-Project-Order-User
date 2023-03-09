package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          string             `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Email       string             `json:"email,omitempty"`
	Password    string             `json:"password,omitempty"`
	Address     string             `json:"address"`
	Orders      []string           `json:"orders"`
	CreatedDate primitive.DateTime `json:"createddate,omitempty"`
	UpdatedDate primitive.DateTime `json:"updateddate"`
}

type Order struct {
	ID      string `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId  string `json:"user_id,omitempty"`
	Product []struct {
		Name     string  `json:"name"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	} `json:"product"`
	Total       float64            `json:"total"`
	CreatedDate primitive.DateTime `json:"createddate,omitempty"`
	UpdatedDate primitive.DateTime `json:"updateddate"`
}
