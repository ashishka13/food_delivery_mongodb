package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"food-delivery/models"
	"food-delivery/services"

	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/mongo"
)

type FoodOrderRequest struct {
	CustomerID     string `json:"customerid" validate:"required"`
	FoodName       string `json:"food_name" validate:"required"`
	Quantity       int    `json:"quantity" validate:"required"`
	RestaurantName string `json:"restaurant_name" validate:"required"`
}

type FoodOrderResponse struct {
	FoodName       string `json:"food_name"`
	Quantity       int    `json:"quantity"`
	RestaurantName string `json:"restaurant_name"`
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

	customer, err := services.CustomerService.GetCustomerByID(ctx, foodOrderRequest.CustomerID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "could not find this customer. First create a customer.", http.StatusForbidden)
			w.Write([]byte("create a customer here: http://localhost:5005/createCustomer"))
			return
		}
		log.Println("error occurred while checkIfCustomerExists", err)
		http.Error(w, "error occurred while getting customer details"+err.Error(), http.StatusInternalServerError)
		return
	}

	foodOrder := convertFoodOrderRequestToFoodOrder(foodOrderRequest)
	createdOrder, err := services.OrderService.CreateOrder(ctx, foodOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("error occurred while creating foodOrderRequest", err)
		return
	}

	foodOrder.ID = createdOrder.ID
	_, err = services.CustomerService.UpdateCustomerOrderID(ctx, customer, foodOrder)
	if err != nil {
		log.Println("error occurred while updating customer order details", err)
		http.Error(w, "error occurred while updating customer order details"+err.Error(), http.StatusInternalServerError)
		return
	}

	displayOrder, err := json.MarshalIndent(FoodOrderResponse{
		FoodName:       foodOrder.FoodName,
		Quantity:       foodOrder.Quantity,
		RestaurantName: foodOrder.RestaurantName}, " ", "  ")
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

func convertFoodOrderRequestToFoodOrder(foodOrderRequest FoodOrderRequest) models.FoodOrder {
	foodOrder := models.FoodOrder{
		CustomerID:     foodOrderRequest.CustomerID,
		FoodName:       foodOrderRequest.FoodName,
		Quantity:       foodOrderRequest.Quantity,
		RestaurantName: foodOrderRequest.RestaurantName,
		PlacedTime:     time.Now(),
	}

	return foodOrder
}
