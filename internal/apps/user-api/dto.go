package user_api

type UserCreateRequest struct {
	Name      string            `json:"name"`
	Email     string            `json:"email" validate:"required,email"`
	Password  string            `json:"password"`
	Addresses []AddressResponse `json:"address"`
}

type UserUpdateRequest struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email" validate:"required,email"`
	Password  string            `json:"password"`
	Addresses []AddressResponse `json:"address"`
}

type UserResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Addresses []AddressResponse `json:"address"`
}

type AddressResponse struct {
	ID       string   `json:"id"`
	Address  string   `json:"address" bson:"address"`
	City     string   `json:"city" bson:"city"`
	District string   `json:"district" bson:"district"`
	Type     []string `json:"type" bson:"type"`
	Default  struct {
		IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
		IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
	} `json:"default" bson:"default"`
}
