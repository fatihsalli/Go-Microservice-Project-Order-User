package user_api

type UserCreateRequest struct {
	Name      string            `json:"name" validate:"required,min=1,max=100"`
	Email     string            `json:"email" validate:"required,email"`
	Password  string            `json:"password" validate:"required,min=1,max=100"`
	Addresses []AddressResponse `json:"address" validate:"required"`
}

type UserUpdateRequest struct {
	ID        string            `json:"id" validate:"required,uuid4"`
	Name      string            `json:"name" validate:"required,min=1,max=100"`
	Email     string            `json:"email" validate:"required,email"`
	Password  string            `json:"password" validate:"required,min=1,max=100"`
	Addresses []AddressResponse `json:"address" validate:"required"`
}

type UserResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Addresses []AddressResponse `json:"address"`
}

type AddressResponse struct {
	ID       string   `json:"id"`
	Address  string   `json:"address" bson:"address" validate:"required"`
	City     string   `json:"city" bson:"city" validate:"required,min=1,max=100"`
	District string   `json:"district" bson:"district" validate:"required,min=1,max=100"`
	Type     []string `json:"type" bson:"type" validate:"required,min=1,max=100"`
	Default  struct {
		IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
		IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
	} `json:"default" bson:"default"`
}
