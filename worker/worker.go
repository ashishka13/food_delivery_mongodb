package worker

import (
	"context"
	"food-delivery/services"
	"food-delivery/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
			// w.CustomerProcess,
			w.DeliveryPersonAssignProcess,
			w.DeliveryPersonPickupProcess,
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
		orders, err := w.Services.OrderService.GetOrderWithFilter(ctx, bson.M{"deliverypersonassigned": false})
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println("error occurred while getting orders", err)
			return
		}

		for _, singleOrder := range orders {
			go func() {
				log.Println("assigning the delivery person for order:", singleOrder.ID.Hex())

				findFilter := bson.M{"id": singleOrder.ID}
				updateSet := bson.M{"$set": bson.M{"deliverypersonassigned": true}}
				err = w.Services.CookingService.UpdateFoodOrderCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the order status", err)
					return
				}

				foundDeliveryPerson, err := w.Services.DeliveryPersonService.GetDeliveryPersonCustom(ctx, bson.M{"busystatus": false})
				if err != nil && err != mongo.ErrNoDocuments {
					log.Print("error occurred while getting free delivery person", err)
					return
				}

				findFilter = bson.M{"id": foundDeliveryPerson.ID}
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

				log.Println("delivery person is on the way to pick up your order from restaurant")

				time.Sleep(time.Second * 10) // time required to reach the restaurant

				updateSet = bson.M{"$set": bson.M{"currentlocation": utils.AtRestaurant}}
				err = w.Services.DeliveryPersonService.UpdateDeliveryPersonCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the delivery person status", err)
					return
				}

				log.Printf("delivery person %s is at restaurant", foundDeliveryPerson.ID.Hex())
			}()
		}
	}
}

func (w *WServices) DeliveryPersonPickupProcess(duhh string) {
	t := time.NewTicker(time.Second * 2)
	for range t.C {

		ctx := context.Background()
		orders, err := w.Services.OrderService.GetOrderWithFilter(ctx, bson.M{"deliverypersonassigned": true, "cookedandready": true})
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println("error occurred while getting orders", err)
			return
		}

		for _, singleOrder := range orders {
			log.Println("code was here")
			go func() {
				foundDeliveryPerson, err := w.Services.DeliveryPersonService.GetDeliveryPersonCustom(ctx, bson.M{"currentorderid": singleOrder.ID.Hex()})
				if err != nil && err != mongo.ErrNoDocuments {
					log.Print("error occurred while getting free delivery person", err)
					return
				}

				findFilter := bson.M{"id": foundDeliveryPerson.ID}
				updateSet := bson.M{"$set": bson.M{"orderpicked": true, "currentlocation": utils.InTransit}}
				err = w.Services.DeliveryPersonService.UpdateDeliveryPersonCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the delivery person status", err)
					return
				}

				time.Sleep(time.Second * 10) // time required to reach the customer

				updateSet = bson.M{"$set": bson.M{"currentlocation": utils.AtCustomer}}
				err = w.Services.DeliveryPersonService.UpdateDeliveryPersonCustom(ctx, findFilter, updateSet)
				if err != nil {
					log.Print("error occurred while updating the delivery person status", err)
					return
				}

				log.Printf("delivery person %s is at customer location", foundDeliveryPerson.ID.Hex())
			}()
		}
	}
}

func (w *WServices) CustomerProcess(duhh string) {

}
