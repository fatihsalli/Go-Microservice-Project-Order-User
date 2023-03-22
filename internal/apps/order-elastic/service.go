package order_elastic

type OrderElasticService struct {
}

func NewOrderElasticService() *OrderElasticService {
	orderService := &OrderElasticService{}
	return orderService
}

type IOrderElasticService interface {
}
