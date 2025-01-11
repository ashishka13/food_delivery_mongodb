package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"food-delivery/models"
	"food-delivery/services"

	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FoodOrderRequest struct {
	CustomerID     string `json:"customerid" validate:"required"`
	FoodName       string `json:"food_name" validate:"required"`
	Quantity       int    `json:"quantity" validate:"required"`
	RestaurantName string `json:"restaurant_name" validate:"required"`
}

type AnyStruct struct {
}

func (a *AnyStruct) postFoodOrder(w http.ResponseWriter, r *http.Request, services services.Services) {
	ctx := r.Context()

	foodOrderRequest := FoodOrderRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&foodOrderRequest)
	if err != nil {
		log.Println("decode foodOrderRequest error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := validator.New()
	err = v.Struct(foodOrderRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("invalid data entered", err)
		return
	}

	customerPresent, err := checkIfCustomerExists(ctx, models.Customer{ID: foodOrderRequest.CustomerID}, services)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("error occurred while getting customer details", err)
		return
	}
	if !customerPresent {
		http.Error(w, "could not find this customer. First create a customer.", http.StatusForbidden)
		w.Write([]byte("create a customer here: http://localhost:5005/createCustomer"))
		return
	}

	foodOrder := convertFoodOrderRequestToFoodOrder(foodOrderRequest)
	_, err = services.OrderService.CreateOrder(ctx, foodOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("error occurred while creating foodOrderRequest", err)
		return
	}

	displayOrder, err := json.MarshalIndent(foodOrderRequest, " ", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("display order marshal error", err)
		return
	}

	_, err = w.Write(displayOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("display order error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func checkIfCustomerExists(ctx context.Context, customer models.Customer, services services.Services) (bool, error) {
	_, err := services.CustomerService.GetCustomer(ctx, customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		log.Println("error occurred while checkIfCustomerExists", err)
		return false, err
	}
	return true, nil
}

func convertFoodOrderRequestToFoodOrder(foodOrderRequest FoodOrderRequest) models.FoodOrder {
	foodOrder := models.FoodOrder{
		ID:             primitive.NewObjectID().Hex(),
		CustomerID:     foodOrderRequest.CustomerID,
		FoodName:       foodOrderRequest.FoodName,
		Quantity:       foodOrderRequest.Quantity,
		RestaurantName: foodOrderRequest.RestaurantName,
		PlacedTime:     time.Now(),
	}

	return foodOrder
}
