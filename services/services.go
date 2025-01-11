package services

import (
	"food-delivery/services/customer"
	deliveryperson "food-delivery/services/delivery-person"
	"food-delivery/services/order"
	"food-delivery/utils"
	"log"
)

type Services struct {
	CustomerService       customer.CustomerServiceInterface
	OrderService          order.OrderServiceInterface
	DeliveryPersonService deliveryperson.DeliveryPersonServicesInterface
}

func InitServices() Services {
	db, err := utils.DatabaseConnect(utils.FoodDelivery)
	if err != nil {
		log.Panic("database connection error", err)
	}

	customerService := customer.NewCustomerServiceInterface(db)
	orderService := order.NewOrderServiceInterface(db)
	deliveryPersonService := deliveryperson.NewDeliveryPersonServicesInterface(db)

	services := Services{
		CustomerService:       customerService,
		OrderService:          orderService,
		DeliveryPersonService: deliveryPersonService,
	}

	return services
}
