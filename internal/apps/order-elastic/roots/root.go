package roots

import (
	"github.com/labstack/gommon/log"
	"sync"
)

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
	group := sync.WaitGroup{}

	group.Add(2)
	go func() {
		defer group.Done()

		if err := o.OrderEventRoot.StartGetOrderAndPushOrder(); err != nil {
			log.Fatalf("OrderEventRoot failed, shutting down the server. | Error: %v\n", err)
		}
	}()
	go func() {
		defer group.Done()

		if err := o.OrderElasticRoot.StartConsumeAndSaveOrder(); err != nil {
			log.Fatalf("OrderElasticRoot failed, shutting down the server. | Error: %v\n", err)
		}
	}()

	group.Wait()
	log.Info("Uygulama sonlandı")
}
