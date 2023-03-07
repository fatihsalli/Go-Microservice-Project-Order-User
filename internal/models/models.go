package models

type User struct {
	ID       string  `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Email    string  `json:"email,omitempty"`
	Password string  `json:"password,omitempty"`
	Address  string  `json:"address,omitempty"`
	Orders   []Order `json:"orders"`
}

type Order struct {
}
