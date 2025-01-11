package models

import (
	"time"
)

type FoodOrder struct {
	ID                     string    `bson:"id" json:"id"`
	CustomerID             string    `bson:"customerid" json:"customerid"`
	FoodName               string    `bson:"foodname" json:"food_name"`
	Quantity               int       `bson:"quantity" json:"quantity"`
	RestaurantName         string    `bson:"restaurantname" json:"restaurant_name"`
	PlacedTime             time.Time `bson:"placedtime" json:"placedtime"`
	CookAssigned           bool      `bson:"cookassigned" json:"cookassigned"`
	DeliveryPersonAssigned bool      `bson:"deliverypersonassigned" json:"deliverypersonassigned"`
	DeliveryPersonID       string    `bson:"deliverypersonid,omitempty" json:"deliverypersonid"`
	CookedAndReady         bool      `bson:"cookedandready" json:"cookedandready"`
	CompleteStatus         bool      `bson:"completestatus" json:"completestatus"`
}

type Restaurant struct {
	ID      string `bson:"id,omitempty" json:"id"`
	Name    string `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
}

type DeliveryPerson struct {
	ID              string `bson:"id,omitempty" json:"id"`
	Name            string `bson:"name" json:"name"`
	BusyStatus      bool   `bson:"busystatus" json:"busystatus"`
	CurrentOrderID  string `bson:"currentorderid" json:"currentorderid"`
	CurrentLocation string `bson:"currentlocation" json:"currentlocation"`
}

type Customer struct {
	ID                 string    `bson:"id" json:"id"`
	Name               string    `bson:"name" json:"name"`
	Address            string    `bson:"address" json:"address"`
	CurrentOrderID     string    `bson:"currentorderid" json:"currentorderid"`
	OrderPlaced        bool      `bson:"orderplaced" json:"orderplaced"`
	PlacedTime         time.Time `bson:"placedtime" json:"placedtime"`
	ReceiveTime        time.Time `bson:"receivetime" json:"receivetime"`
	DeliveryPersonName string    `bson:"deliverypersonname" json:"deliverypersonname"`
	InProcess          bool      `bson:"inprocess" json:"inprocess"`
}
