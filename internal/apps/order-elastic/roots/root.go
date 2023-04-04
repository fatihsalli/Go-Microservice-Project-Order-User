package roots

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
			o.OrderEventRoot.Logger.Fatalf("Order sync service (StartGetOrderAndPushOrder) failed, shutting down the server. | Error: %v\n", err)
		}
	}()
	if err := o.OrderElasticRoot.StartConsumeAndSaveOrder(); err != nil {
		o.OrderElasticRoot.Logger.Fatalf("Order sync service (StartConsumeAndSaveOrder) failed, shutting down the server. | Error: %v\n", err)
	}
}
