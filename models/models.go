package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FoodOrder struct {
	ID             string             `bson:"id" json:"id"`
	FoodName       string             `bson:"foodname" json:"food_name"`
	Quantity       int                `bson:"quantity" json:"quantity"`
	RestaurantName string             `bson:"restaurantname" json:"restaurant_name"`
	Address        string             `bson:"address" json:"address"`
	PlacedTime     time.Time          `bson:"placedtime" json:"placed_time"`
	CookAssigned   bool               `bson:"cookassigned" json:"cookassigned"`
	BoyAssigned    bool               `bson:"boyassigned" json:"boyassigned"`
	BoyID          primitive.ObjectID `bson:"boyid,omitempty" json:"boy_id"`
	CookedAndReady bool               `bson:"cookedandready" json:"cookedandready"`
	CompleteStatus bool               `bson:"completestatus" json:"complete_status"`
}

type Restaurant struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Address string             `bson:"address" json:"address"`
}

type DeliveryBoy struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name"`
	BusyStatus      bool               `bson:"busystatus" json:"busystatus"`
	CurrentOrderID  string             `bson:"currentorderid" json:"current_order_id"`
	CurrentLocation string             `bson:"currentlocation" json:"current_location"`
}

type Customer struct {
	ID             string    `bson:"id" json:"id"`
	Name           string    `bson:"name" json:"name"`
	Address        string    `bson:"address" json:"address"`
	CurrentOrderID string    `bson:"currentorderid" json:"current_order_id"`
	OrderPlaced    bool      `bson:"orderplaced" json:"order_placed"`
	PlacedTime     time.Time `bson:"placedtime" json:"placed_time"`
	ReceiveTime    time.Time `bson:"receivetime" json:"receive_time"`
	BoyName        string    `bson:"boyname" json:"boy_name"`
	InProcess      bool      `bson:"inprocess" json:"inprocess"`
}
