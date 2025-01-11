package controller

import (
	"food-delivery/services"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	Route *mux.Router
}

func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request, services.Services)) *mux.Route {
	return r.Route.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		services := services.InitServices()
		f(w, r, services)
	})
}
