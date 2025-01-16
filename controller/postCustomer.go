package controller

import (
	"encoding/json"
	"food-delivery/models"
	"food-delivery/services"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerCreateRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
}

type CustomerPostResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func postCustomer(w http.ResponseWriter, r *http.Request, services services.Services) {
	ctx := r.Context()
	customerCreateRequest := CustomerCreateRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&customerCreateRequest)
	if err != nil {
		log.Println("decode customerCreateRequest error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := validator.New()
	err = v.Struct(customerCreateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("invalid data entered", err)
		return
	}

	customerFound, err := services.CustomerService.GetCustomerByPhone(ctx, customerCreateRequest.Phone)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("error occurred while checkIfCustomerExists", err)
		http.Error(w, "error occurred while matching customer details", http.StatusInternalServerError)
		return
	}
	if customerFound != nil {
		log.Println("customer already present")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("customer already present with these details"))
		return
	}

	customer := convertCustomerCreateRequestToCustomer(customerCreateRequest)
	createdCustomer, err := services.CustomerService.CreateCustomer(ctx, customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("error occurred while creating customerCreateRequest", err)
		return
	}

	displayCustomer, err := json.MarshalIndent(CustomerPostResponse{
		ID:      createdCustomer.ID.Hex(),
		Name:    customer.Name,
		Address: customer.Address,
		Phone:   customer.Phone}, " ", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("display customer marshal error", err)
		return
	}

	_, err = w.Write(displayCustomer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("display customer error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func convertCustomerCreateRequestToCustomer(customerCreateRequest CustomerCreateRequest) models.Customer {
	customer := models.Customer{
		Name:    customerCreateRequest.Name,
		Address: customerCreateRequest.Address,
		Phone:   customerCreateRequest.Phone,
		Orders:  []models.CustomerOrders{},
	}

	return customer
}
