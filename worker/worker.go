package worker

import (
	"context"
	"encoding/json"
	"food-delivery/models"
	"food-delivery/services"
	"food-delivery/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func MainWorker() {
	log.Println("inside Mainworker")
	RunProcesses()
}

type WServices struct {
	Services services.Services
}

type JobRunners struct {
	Runners []func(duhh string)
}

func RunProcesses() {
	services := services.InitServices()
	w := WServices{Services: services}

	j := JobRunners{
		Runners: []func(duhh string){
			w.CustomerProcess,
			w.DeliveryPersonAssignProcess,
			w.RestaurantHandoverOrderProcesss,
			w.RestaurantCookingProcess,
		},
	}

	for _, singleJob := range j.Runners {
		go singleJob("magic duhh")
	}
}

func (w *WServices) RestaurantCookingProcess(duhh string) {
	t := time.NewTicker(time.Second * 2)
	for range t.C {

		ctx := context.Background()
		orders, err := w.Services.OrderService.GetOrderWithFilter(ctx, bson.M{"cookassigned": false})
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println("error occurred while getting orders", err)
			return
		}

		for _, singleOrder := range orders {
			go func() {
				log.Println("cooking the order:", singleOrder.ID.Hex())

				findFilter := bson.M{"id": singleOrder.ID}
				updateSet := bson.M{"$set": bson.M{"cookassigned": true}}
				err = w.Services.CookingService.UpdateFoodOrderCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the order status", err)
					return
				}

				time.Sleep(time.Second * 10) // time required for cooking the order

				updateSet = bson.M{"$set": bson.M{"cookedandready": true}}
				err = w.Services.CookingService.UpdateFoodOrderCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the order status", err)
					return
				}

				log.Printf("order %s is ready", singleOrder.ID.Hex())
			}()
		}
	}
}

func (w *WServices) DeliveryPersonAssignProcess(duhh string) {
	t := time.NewTicker(time.Second * 2)
	for range t.C {
		ctx := context.Background()

		ordersAndPersons, err := w.clubDeliveryPersonAndOrder(ctx)
		if err != nil {
			log.Println("error occurred while clubDeliveryPersonAndOrder", err)
			return
		}

		for singlePerson, singleOrder := range ordersAndPersons {
			go func() {
				log.Println("assigning the delivery person for order:", singleOrder.ID.Hex())

				findFilter := bson.M{"id": singleOrder.ID}
				updateSet := bson.M{"$set": bson.M{"deliverypersonassigned": true}}
				err = w.Services.CookingService.UpdateFoodOrderCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the order status", err)
					return
				}

				findFilter = bson.M{"id": singlePerson.ID}
				updateSet = bson.M{"$set": bson.M{
					"currentorderid":    singleOrder.ID.Hex(),
					"currentcustomerid": singleOrder.CustomerID,
					"currentlocation":   utils.InTransit,
					"busystatus":        true}}
				err = w.Services.DeliveryPersonService.UpdateDeliveryPersonCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the delivery person status", err)
					return
				}

				customerObjectID, _ := primitive.ObjectIDFromHex(singleOrder.CustomerID)
				findFilter = bson.M{"id": customerObjectID, "orders.currentorderid": singleOrder.ID.Hex()}
				updateSet = bson.M{"$set": bson.M{"orders.$.deliverypersonid": singlePerson.ID.Hex()}}

				findQurty, _ := json.Marshal(findFilter)
				updateq, _ := json.Marshal(updateSet)

				log.Printf("find %s, update %s", findQurty, updateq)

				err = w.Services.CustomerService.UpdateCustomerOrderCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the customer delivery person details", err)
					return
				}

				log.Println("delivery person is on the way to pick up your order from restaurant")

				time.Sleep(time.Second * 10) // time required to reach the restaurant

				findFilter = bson.M{"id": singlePerson.ID}
				updateSet = bson.M{"$set": bson.M{"currentlocation": utils.AtRestaurant}}
				err = w.Services.DeliveryPersonService.UpdateDeliveryPersonCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the delivery person status", err)
					return
				}

				log.Printf("delivery person %s is at restaurant", singlePerson.ID.Hex())
			}()
		}
	}
}

func (w *WServices) RestaurantHandoverOrderProcesss(duhh string) {
	t := time.NewTicker(time.Second * 2)
	for range t.C {
		ctx := context.Background()

		ordersAndPersons, err := w.clubReadyOrdersAndDeliveryPerson(ctx)
		if err != nil {
			log.Println("error occurred while clubReadyOrdersAndDeliveryPerson", err)
			return
		}

		for singleOrder, singlePerson := range ordersAndPersons {
			go func() {
				findFilter := bson.M{"id": singleOrder.ID}
				updateSet := bson.M{"$set": bson.M{"pickedstatus": true}}
				err = w.Services.OrderService.UpdateOrder(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the order pick status", err)
					return
				}

				findFilter = bson.M{"id": singlePerson.ID}
				updateSet = bson.M{"$set": bson.M{"orderpicked": true, "currentlocation": utils.InTransit}}
				err = w.Services.DeliveryPersonService.UpdateDeliveryPersonCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the delivery person status", err)
					return
				}

				log.Printf("order picked. Delivery person %s is on the way", singlePerson.ID.Hex())

				time.Sleep(time.Second * 10) // time required to reach the customer

				updateSet = bson.M{"$set": bson.M{"currentlocation": utils.AtCustomer}}
				err = w.Services.DeliveryPersonService.UpdateDeliveryPersonCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the delivery person status", err)
					return
				}

				log.Printf("delivery person %s is at customer location", singlePerson.ID.Hex())
			}()
		}
	}
}

func (w *WServices) CustomerProcess(duhh string) {
	t := time.NewTicker(time.Second * 2)
	for range t.C {
		ctx := context.Background()
		customerAndReadyPickOrders, err := w.clubWaitingCustomerAndPickReadyOrders(ctx)
		if err != nil {
			log.Println("error occurred while clubWaitingCustomerAndPickReadyOrders", err)
			return
		}

		for singleCustomer, singleOrder := range customerAndReadyPickOrders {
			log.Printf("customer %s is getting order from delivery person %s", singleCustomer.ID.Hex(), singleOrder.DeliveryPersonID)
			time.Sleep(time.Second * 7) //time taken by customer to receive the order

			findFilter := bson.M{"id": singleCustomer.ID, "orders.currentorderid": singleOrder.CurrentOrderID}
			updateSet := bson.M{"$set": bson.M{"orders.$.receivetime": time.Now()}}
			if err = w.Services.CustomerService.UpdateCustomerOrderCustom(ctx, findFilter, updateSet); err != nil {
				log.Println(err)
				return
			}
			log.Printf("order %s picked up by customer", singleOrder.CurrentOrderID)

			deliveryObjectID, _ := primitive.ObjectIDFromHex(singleOrder.DeliveryPersonID)
			findFilter = bson.M{"id": deliveryObjectID, "currentlocation": utils.AtCustomer}
			updateSet = bson.M{"$set": bson.M{"busystatus": false, "currentlocation": utils.Dock, "currentorderid": "", "currentcustomerid": "", "orderpicked": false}}
			if err = w.Services.DeliveryPersonService.UpdateDeliveryPersonCustom(ctx, findFilter, updateSet); err != nil {
				log.Println("error occurred while updating delivery person", err)
				return
			}
		}
	}
}

func (w *WServices) clubDeliveryPersonAndOrder(ctx context.Context) (map[models.DeliveryPerson]models.FoodOrder, error) {
	findFilter := bson.M{"deliverypersonassigned": false}
	orders, err := w.Services.OrderService.GetOrderWithFilter(ctx, findFilter)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("error occurred while getting orders", err)
		return nil, err
	}

	findFilter = bson.M{"busystatus": false}
	deliveryPersons, err := w.Services.DeliveryPersonService.GetDeliveryPersonsCustom(ctx, findFilter)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("error occurred while getting delivery persons", err)
		return nil, err
	}

	ordersAndPersons := make(map[models.DeliveryPerson]models.FoodOrder)

	lengthLimit := len(deliveryPersons)
	if lengthLimit > len(orders) {
		lengthLimit = len(orders)
	}

	for i := 0; i < lengthLimit; i++ {
		ordersAndPersons[deliveryPersons[i]] = orders[i]
	}

	return ordersAndPersons, nil
}

func (w *WServices) clubReadyOrdersAndDeliveryPerson(ctx context.Context) (map[models.FoodOrder]models.DeliveryPerson, error) {
	ordersAndPersons := make(map[models.FoodOrder]models.DeliveryPerson)

	orders, err := w.Services.OrderService.GetOrderWithFilter(ctx, bson.M{"deliverypersonassigned": true, "cookedandready": true, "pickedstatus": false})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("error occurred while getting orders", err)
		return nil, err
	}

	for _, singleOrder := range orders {
		findFilter := bson.M{"currentorderid": singleOrder.ID.Hex(), "currentlocation": utils.AtRestaurant}
		foundDeliveryPerson, err := w.Services.DeliveryPersonService.GetDeliveryPersonCustom(ctx, findFilter)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Print("error occurred while getting free delivery person", err)
			return nil, err
		}

		ordersAndPersons[singleOrder] = foundDeliveryPerson
	}

	return ordersAndPersons, nil
}

func (w *WServices) clubWaitingCustomerAndPickReadyOrders(ctx context.Context) (map[*models.Customer]models.CustomerOrders, error) {
	customerAndReadyPickOrders := make(map[*models.Customer]models.CustomerOrders)

	findFilter := bson.M{"orders": bson.M{"$elemMatch": bson.M{"receivetime": time.Time{}}}}
	customers, err := w.Services.CustomerService.GetAllWaitingCustomersCustom(ctx, findFilter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Println("error occurred while getting customers list", err)
		return nil, err
	}

	for _, singleCustomer := range customers {
		findFilter = bson.M{"currentcustomerid": singleCustomer.ID.Hex(), "orderpicked": true, "currentlocation": utils.AtCustomer}

		deliverypersons, err := w.Services.DeliveryPersonService.GetDeliveryPersonsCustom(ctx, findFilter)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}
			log.Println("error occurred while getting delivery persons", err)
			return nil, err
		}

		for _, singleDeliveryPerson := range deliverypersons {
			customerAndReadyPickOrders[&singleCustomer] = models.CustomerOrders{CurrentOrderID: singleDeliveryPerson.CurrentOrderID, DeliveryPersonID: singleDeliveryPerson.ID.Hex()}
		}
	}
	return customerAndReadyPickOrders, nil
}
