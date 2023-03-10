package user_api

type UserCreateRequest struct {
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Address  []AddressResponse `json:"address"`
}

type UserUpdateRequest struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Address  []AddressResponse `json:"address"`
}

type UserResponse struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Email   string            `json:"email"`
	Address []AddressResponse `json:"address"`
}

type AddressResponse struct {
	Address  string   `json:"address"`
	City     string   `json:"city"`
	District string   `json:"district"`
	Type     []string `json:"type"`
}
