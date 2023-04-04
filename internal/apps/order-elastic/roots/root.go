package roots

import "github.com/labstack/gommon/log"

type OrderSyncService struct {
	OrderElasticRoot *OrderElasticRoot
	OrderEventRoot   *OrderEventRoot
}

func NewOrderSyncService(orderElasticRoot *OrderElasticRoot, orderEventRoot *OrderEventRoot) *OrderSyncService {
	return &OrderSyncService{
		OrderElasticRoot: orderElasticRoot,
		OrderEventRoot:   orderEventRoot,
	}
}

func (o OrderSyncService) Start() {
	go func() {
		if err := o.OrderEventRoot.StartGetOrderAndPushOrder(); err != nil {
			log.Fatalf("OrderEventRoot failed, shutting down the server. | Error: %v\n", err)
		}
	}()
	if err := o.OrderElasticRoot.StartConsumeAndSaveOrder(); err != nil {
		log.Fatalf("OrderElasticRoot failed, shutting down the server. | Error: %v\n", err)
	}
}
