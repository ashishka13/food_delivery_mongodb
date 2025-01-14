package controller

import (
	"context"
	"log"
	"net/http"

	"food-delivery/models"
	"food-delivery/services"
	deliveryperson "food-delivery/services/delivery-person"
	"food-delivery/utils"

	"github.com/gorilla/mux"
)

// MyController ...
func MyController() {
	router := mux.NewRouter()
	r := Router{Route: router}
	d := AnyStruct{} // this way we can create a method on struct and pass parameters in our custom struct
	// we can use them inside our API function as well.

	a, _ := utils.DatabaseConnect(utils.FoodDelivery)
	boysvc := services.Services{
		DeliveryPersonService: deliveryperson.NewDeliveryPersonServicesInterface(a),
	}
	boysvc.DeliveryPersonService.CreateDeliveryPerson(context.Background(), models.DeliveryPerson{})

	r.HandleFunc("/", welcome).Methods("GET", "PUT", "POST", "DELETE")
	r.HandleFunc("/orderFood", d.postFoodOrder).Methods("POST")
	r.HandleFunc("/createCustomer", postCustomer).Methods("POST")
	r.HandleFunc("/getCustomer", getCustomer).Methods("GET")

	log.Println("listening..")
	http.ListenAndServe(":5005", r.Route)
}

func welcome(w http.ResponseWriter, r *http.Request, services services.Services) {
	w.Write([]byte("welcome to food delivery service, choose restaurant and place order"))
}
