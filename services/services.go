package services

import (
	"food-delivery/services/cooking"
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
	CookingService        cooking.CookingServiceInterface
}

func InitServices() Services {
	db, err := utils.DatabaseConnect(utils.FoodDelivery)
	if err != nil {
		log.Panic("database connection error", err)
	}

	customerService := customer.NewCustomerServiceInterface(db)
	orderService := order.NewOrderServiceInterface(db)
	deliveryPersonService := deliveryperson.NewDeliveryPersonServicesInterface(db)
	cookingService := cooking.NewCookingServiceInterface(db)

	services := Services{
		CustomerService:       customerService,
		OrderService:          orderService,
		DeliveryPersonService: deliveryPersonService,
		CookingService:        cookingService,
	}

	return services
}
