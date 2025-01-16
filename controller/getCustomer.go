package controller

import (
	"encoding/json"
	"food-delivery/services"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerGetRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type CustomerGetResponse struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Address string           `json:"address"`
	Phone   string           `json:"phone"`
	Orders  []CustomerOrders `json:"orders"`
}

type CustomerOrders struct {
	CurrentOrderID string `json:"currentorderid"`
	OrderName      string `json:"ordername"`
}

func getCustomer(w http.ResponseWriter, r *http.Request, services services.Services) {
	ctx := r.Context()
	CustomerGetRequest := CustomerGetRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&CustomerGetRequest)
	if err != nil {
		log.Println("decode CustomerGetRequest error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := validator.New()
	err = v.Struct(CustomerGetRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("invalid data entered", err)
		return
	}

	customerFound, err := services.CustomerService.GetCustomerByPhone(ctx, CustomerGetRequest.Phone)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "could not find this customer.", http.StatusNotFound)
			return
		}
		log.Println("error occurred while checkIfCustomerExists", err)
		http.Error(w, "error occurred while matching customer details", http.StatusInternalServerError)
		return
	}

	customerOrders := make([]CustomerOrders, len(customerFound.Orders))
	for key, val := range customerFound.Orders {
		customerOrders[key].CurrentOrderID = val.CurrentOrderID
		customerOrders[key].OrderName = val.OrderName
	}

	displayCustomer, err := json.MarshalIndent(CustomerGetResponse{
		ID:      customerFound.ID.Hex(),
		Name:    customerFound.Name,
		Address: customerFound.Address,
		Phone:   customerFound.Phone,
		Orders:  customerOrders}, " ", "  ")

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
