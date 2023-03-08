package order_api

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

type OrderResponse struct {
	ID      string `json:"id"`
	UserId  string `json:"user_id"`
	Product []struct {
		Name     string  `json:"name"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
	Total float64 `json:"total"`
}
