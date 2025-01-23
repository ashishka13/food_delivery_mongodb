package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FoodOrder struct {
	ID                     primitive.ObjectID `bson:"id" json:"id"`
	CustomerID             string             `bson:"customerid" json:"customerid"`
	FoodName               string             `bson:"foodname" json:"food_name"`
	Quantity               int                `bson:"quantity" json:"quantity"`
	RestaurantName         string             `bson:"restaurantname" json:"restaurant_name"`
	PlacedTime             time.Time          `bson:"placedtime" json:"placedtime"`
	CookAssigned           bool               `bson:"cookassigned" json:"cookassigned"`
	DeliveryPersonAssigned bool               `bson:"deliverypersonassigned" json:"deliverypersonassigned"`
	DeliveryPersonID       string             `bson:"deliverypersonid,omitempty" json:"deliverypersonid"`
	CookedAndReady         bool               `bson:"cookedandready" json:"cookedandready"`
	PickedStatus           bool               `bson:"pickedstatus" json:"pickedstatus"`
}

type Restaurant struct {
	ID      primitive.ObjectID `bson:"id" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Address string             `bson:"address" json:"address"`
}

type DeliveryPerson struct {
	ID                primitive.ObjectID `bson:"id" json:"id"`
	Name              string             `bson:"name" json:"name"`
	BusyStatus        bool               `bson:"busystatus" json:"busystatus"`
	CurrentOrderID    string             `bson:"currentorderid" json:"currentorderid"`
	CurrentCustomerID string             `bson:"currentcustomerid" json:"currentcustomerid"`
	CurrentLocation   string             `bson:"currentlocation" json:"currentlocation"`
	OrderPicked       bool               `bson:"orderpicked" json:"orderpicked"`
}

type Customer struct {
	ID      primitive.ObjectID `bson:"id" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Address string             `bson:"address" json:"address"`
	Phone   string             `bson:"phone" json:"phone"`
	Orders  []CustomerOrders   `bson:"orders" json:"orders"`
}

type CustomerOrders struct {
	OrderName        string    `bson:"ordername" json:"ordername"`
	CurrentOrderID   string    `bson:"currentorderid" json:"currentorderid"`
	PlacedTime       time.Time `bson:"placedtime" json:"placedtime"`
	ReceiveTime      time.Time `bson:"receivetime" json:"receivetime"`
	DeliveryPersonID string    `bson:"deliverypersonid,omitempty" json:"deliverypersonid"`
}
