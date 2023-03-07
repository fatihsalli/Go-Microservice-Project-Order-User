package app

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type UserUpdateRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type OrderCreateRequest struct {
	UserId  string `json:"user_id"`
	Product []struct {
		Name     string  `json:"name"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
	Total float64 `json:"total"`
}

type OrderUpdateRequest struct {
	ID      string `json:"id"`
	UserId  string `json:"user_id"`
	Product []struct {
		Name     string  `json:"name"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
	Total float64 `json:"total"`
}
