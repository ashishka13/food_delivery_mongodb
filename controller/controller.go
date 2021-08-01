package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"foodDelivery/models"
	"foodDelivery/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MyController ...
func MyController() {
	r := mux.NewRouter()

	r.HandleFunc("/", welcome).Methods("GET", "PUT", "POST", "DELETE")
	r.HandleFunc("/findFood", findRestaurant).Methods("GET")
	r.HandleFunc("/orderFood/{customerName}", postOrder).Methods("POST")

	log.Println("listening..")
	http.ListenAndServe(":5005", r)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to food delivery service, choose restaurant and place order"))
}

func findRestaurant(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to findFood"))
}

func postOrder(w http.ResponseWriter, r *http.Request) {
	db := utils.DatabaseConnect(utils.FoodDelivery)
	ctx := r.Context()

	params := mux.Vars(r)

	customerName := params["customerName"]
	singleOrder := models.FoodOrder{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&singleOrder)
	if err != nil {
		log.Println("decode food order error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpstatus, err := createOrder(ctx, db, singleOrder, customerName)
	log.Println("status", httpstatus)
	if err != nil {
		http.Error(w, err.Error(), httpstatus)
		fmt.Fprintf(w, "%v ", httpstatus)
		log.Println("insert order error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func createOrder(ctx context.Context, db *mongo.Database, singleOrder models.FoodOrder, customerName string) (httpstatus int, err error) {
	prevCustomer := models.Customer{}
	newCustomer := models.Customer{}
	newInsertID := primitive.NewObjectID().Hex()
	oldcustomer := true
	err = db.Collection(utils.Customers).FindOne(context.Background(), bson.M{"name": customerName}).Decode(&prevCustomer)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("createOrder customer find one errors", prevCustomer)
		httpstatus = http.StatusInternalServerError
		return
	}
	if err == mongo.ErrNoDocuments {
		log.Println("customer not found, creating new customer")
		newCustomer = models.Customer{
			ID:      newInsertID,
			Name:    customerName,
			Address: "temporary address update later",
		}
		_, err2 := db.Collection(utils.Customers).InsertOne(ctx, newCustomer)
		if err2 != nil {
			log.Println("new customer InsertOne error occurred", err2)
			httpstatus = http.StatusInternalServerError
			return
		}
		oldcustomer = false
	}
	if prevCustomer.OrderPlaced {
		log.Println("this customer order already in progress", prevCustomer.Name)
		httpstatus = http.StatusConflict
		err = errors.New("customer already present")
		return
	}

	singleOrder.ID = primitive.NewObjectID().Hex()
	placedtime := time.Now()
	singleOrder.PlacedTime = placedtime

	_, err = db.Collection(utils.Orders).InsertOne(ctx, singleOrder)
	if err != nil {
		log.Println("orders InsertOne error occurred", err)
		httpstatus = http.StatusInternalServerError
		return
	}
	cfilter := bson.M{"id": prevCustomer.ID}
	if !oldcustomer {
		cfilter = bson.M{"id": newInsertID}
	}
	cupdate := bson.M{"$set": bson.M{"currentorderid": singleOrder.ID, "orderplaced": true, "placedtime": placedtime}}
	res, err := db.Collection(utils.Customers).UpdateOne(ctx, cfilter, cupdate)
	if err != nil || res.ModifiedCount == 0 {
		log.Println("find single customer error", err)
		return
	}
	return
}
