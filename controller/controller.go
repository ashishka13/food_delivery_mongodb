package controller

import (
	"log"
	"net/http"

	"food-delivery/services"

	"github.com/gorilla/mux"
)

// MyController ...
func MyController() {
	router := mux.NewRouter()
	r := Router{Route: router}
	d := AnyStruct{} // this way we can create a method on struct and pass parameters in our custom struct
	// we can use them inside our API function as well.

	r.HandleFunc("/", welcome).Methods("GET", "PUT", "POST", "DELETE")
	r.HandleFunc("/orderFood", d.postFoodOrder).Methods("POST")
	r.HandleFunc("/createCustomer", welcome).Methods("POST")

	log.Println("listening..")
	http.ListenAndServe(":5005", r.Route)
}

func welcome(w http.ResponseWriter, r *http.Request, services services.Services) {
	w.Write([]byte("welcome to food delivery service, choose restaurant and place order"))
}
