package processor

import (
	"context"
	"encoding/json"
	"foodDelivery/models"
	"foodDelivery/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func MainProcessor() {
	log.Println("inside MainProcessor")
	go CheckOrdersRestaurant()
	go CheckOrdersDeliveryMan()
	go CheckPendingProcessOrders()
	// select {}
}

func CheckOrdersRestaurant() {
	log.Println("inside CheckOrdersRestaurant")
	db := utils.DatabaseConnect(utils.FoodDelivery)
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		resultset, err := db.Collection(utils.Orders).Find(context.Background(), bson.M{"cookassigned": false})
		if err != nil {
			log.Println("CheckOrdersRestaurant Orders find all ", err)
			return
		}

		for resultset.Next(context.Background()) {
			singleFood := models.FoodOrder{}
			err = resultset.Decode(&singleFood)
			if err != nil {
				log.Println("decode error")
				return
			}
			log.Println("singleFood CheckOrdersRestaurant", singleFood)
			StartCookFood(db, singleFood)
		}
		resultset.Close(context.Background())
	}
}

func CheckOrdersDeliveryMan() {
	log.Println("inside CheckOrdersDeliveryMan")
	db := utils.DatabaseConnect(utils.FoodDelivery)
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		resultset, err := db.Collection(utils.Orders).Find(context.Background(), bson.M{"boyassigned": false, "cookassigned": true})
		if err != nil {
			log.Println("CheckOrdersDeliveryMan Orders find all", err)
			return
		}
		for resultset.Next(context.Background()) {
			singleFood := models.FoodOrder{}
			err = resultset.Decode(&singleFood)
			if err != nil {
				log.Println("decode error")
				return
			}

			dfilter := primitive.M{"currentorderid": singleFood.ID}
			count, err := db.Collection(utils.DeliveryBoys).CountDocuments(context.Background(), dfilter)
			if err != nil {
				log.Println("TowardsRestaurant count error DeliveryBoys", err)
				continue
			}
			if count > 0 {
				log.Println("this order has already a boy assigned")
				manageTimes(db, singleFood)
				continue
			}

			singleBoy := models.DeliveryBoy{}
			err = db.Collection(utils.DeliveryBoys).FindOne(context.Background(), bson.M{"busystatus": false}).Decode(&singleBoy)
			if err != nil {
				log.Println("CheckOrdersDeliveryMan DeliveryBoys findone", err, singleBoy.Name)
				continue
			}
			log.Println("singleBoy CheckOrdersDeliveryMan", singleBoy.Name)
			TowardsRestaurant(db, singleBoy, singleFood)
		}
		resultset.Close(context.Background())
	}
}

func CheckPendingProcessOrders() {
	log.Println("inside CheckPendingProcessOrders")

	db := utils.DatabaseConnect(utils.FoodDelivery)
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		resultSet, err := db.Collection(utils.Orders).Find(context.Background(), bson.M{"cookassigned": true, "boyassigned": true, "completestatus": false})
		if err != nil {
			log.Println("CheckPendingProcessOrders orders find error", err)
			continue
		}
		for resultSet.Next(context.Background()) {
			log.Println("code was here")
			singleOrder := models.FoodOrder{}
			err = resultSet.Decode(&singleOrder)
			if err != nil {
				log.Println("decode error", err)
				continue
			}
			manageTimes(db, singleOrder)
		}
		resultSet.Close(context.Background())
	}
}

func StartCookFood(db *mongo.Database, singleFood models.FoodOrder) {
	log.Println("inside CookFood")

	ofilter := bson.M{"id": singleFood.ID}
	oupdate := bson.M{"$set": bson.M{"cookassigned": true}}
	_, err := db.Collection(utils.Orders).UpdateOne(context.Background(), ofilter, oupdate)
	if err != nil {
		log.Println("StartCookFood Orders update error occurred", err)
		return
	}

	cfilter := bson.M{"currentorderid": singleFood.ID}
	cupdate := bson.M{"$set": bson.M{"inprocess": true}}
	_, err = db.Collection(utils.Customers).UpdateOne(context.Background(), cfilter, cupdate)
	if err != nil {
		log.Println("StartCookFood Orders update error occurred", err)
		return
	}

	for i := 0; i < singleFood.Quantity; i++ {
		log.Println("cooking the ", singleFood.FoodName)
		time.Sleep(time.Second * 3)
	}

	oupdate = bson.M{"$set": bson.M{"cookedandready": true}}
	_, err = db.Collection(utils.Orders).UpdateOne(context.Background(), ofilter, oupdate)
	if err != nil {
		log.Println("Orders boy update error occurred", err)
		return
	}
}

func TowardsRestaurant(db *mongo.Database, singleBoy models.DeliveryBoy, singleFood models.FoodOrder) {
	log.Println("inside TowardsRestaurant")

	dfilter := primitive.M{"_id": singleBoy.ID}
	dupdate := primitive.M{"$set": primitive.M{"currentlocation": utils.InTransit, "busystatus": true, "currentorderid": singleFood.ID}}
	updateresult, err := db.Collection(utils.DeliveryBoys).UpdateOne(context.Background(), dfilter, dupdate)
	stringfilter, _ := json.Marshal(dfilter)
	log.Printf("%s", stringfilter)
	if err != nil || updateresult.MatchedCount == 0 {
		log.Println("TowardsRestaurant UpdateOne error DeliveryBoys1", err)
		return
	}

	cfilter := primitive.M{"currentorderid": singleFood.ID}
	cupdate := primitive.M{"$set": primitive.M{"boyname": singleBoy.Name}}
	updateresult, err = db.Collection(utils.Customers).UpdateOne(context.Background(), cfilter, cupdate)
	if err != nil || updateresult.ModifiedCount == 0 {
		log.Println("TowardsRestaurant UpdateOne error utils.Customers", err)
		return
	}

	ofilter := primitive.M{"id": singleFood.ID}
	oupdate := primitive.M{"$set": primitive.M{"boyassigned": true}}
	_, err = db.Collection(utils.Orders).UpdateOne(context.Background(), ofilter, oupdate)
	if err != nil {
		log.Println("TowardsRestaurant UpdateOne error utils.Orders", err)
		return
	}

	for i := 0; i < singleFood.Quantity; i++ {
		log.Print(".")
		time.Sleep(time.Second * 3)
	}

	dupdate = primitive.M{"$set": primitive.M{"currentlocation": utils.AtRestaurant}}
	_, err = db.Collection(utils.DeliveryBoys).UpdateOne(context.Background(), dfilter, dupdate)
	log.Println("TowardsRestaurant UpdateOne error DeliveryBoys2", err)
}

func WaitToPickFromRestaurant(db *mongo.Database, singleBoy models.DeliveryBoy) {
	log.Println("inside WaitToPickFromRestaurant")

	for i := 0; i < 3; i++ {
		log.Print(".")
		time.Sleep(time.Second * 3)
	}

	dfilter := primitive.M{"_id": singleBoy.ID}
	dupdate := primitive.M{"$set": primitive.M{"currentlocation": utils.InTransit}}
	_, err := db.Collection(utils.DeliveryBoys).UpdateOne(context.Background(), dfilter, dupdate)
	log.Println("WaitToPickFromRestaurant UpdateOne error", err)
}

func TowardsCustomer(db *mongo.Database, singleBoy models.DeliveryBoy) {
	log.Println("inside TowardsCustomer")

	for i := 0; i < 3; i++ {
		log.Print(".")
		time.Sleep(time.Second * 3)
	}

	dfilter := primitive.M{"_id": singleBoy.ID}
	dupdate := primitive.M{"$set": primitive.M{"currentlocation": utils.AtCustomer}}
	_, err := db.Collection(utils.DeliveryBoys).UpdateOne(context.Background(), dfilter, dupdate)
	if err != nil {
		log.Println("TowardsRestaurant UpdateOne error", err)
	}
}

func WaitToGiveCustomer(db *mongo.Database, singleFood models.FoodOrder) (err error) {
	log.Println("inside WaitToGiveCustomer")
	for i := 0; i < 3; i++ {
		log.Print(".")
		time.Sleep(time.Second * 5)
	}

	cfilter := primitive.M{"currentorderid": singleFood.ID}
	cupdate := primitive.M{"$set": primitive.M{"receivetime": time.Now(), "currentorderid": "", "orderplaced": false, "inprocess": false}}
	_, err = db.Collection(utils.Customers).UpdateOne(context.Background(), cfilter, cupdate)
	if err != nil {
		log.Println("WaitToGiveCustomer Customers UpdateOne error", err)
		return
	}
	ofilter := primitive.M{"id": singleFood.ID}
	oupdate := primitive.M{"$set": primitive.M{"completestatus": true}}
	_, err = db.Collection(utils.Orders).UpdateOne(context.Background(), ofilter, oupdate)
	if err != nil {
		log.Println("WaitToGiveCustomer Orders UpdateOne error", err)
		return
	}
	dfilter := primitive.M{"currentorderid": singleFood.ID}
	dupdate := primitive.M{"$set": primitive.M{"busystatus": false, "currentorderid": "", "currentlocation": utils.Dock}}
	_, err = db.Collection(utils.DeliveryBoys).UpdateOne(context.Background(), dfilter, dupdate)
	if err != nil {
		log.Println("WaitToGiveCustomer DeliveryBoys UpdateOne error", err)
		return
	}
	return
}

func manageTimes(db *mongo.Database, singleFood models.FoodOrder) {
	cfilter := primitive.M{"orderplaced": true}
	filterstring, err := json.Marshal(&cfilter)
	log.Printf("filter %s error %v ", filterstring, err)
	count, err := db.Collection(utils.Customers).CountDocuments(context.Background(), cfilter)
	if err != nil {
		log.Println("WaitToGiveCustomer Customers count orders error", err)
		return
	}
	log.Println("count ", count)
	if count == 0 {
		log.Println("no pending orders")
		return
	}
	singleBoy := models.DeliveryBoy{}
	err = db.Collection(utils.DeliveryBoys).FindOne(context.Background(), bson.M{"currentorderid": singleFood.ID}).Decode(&singleBoy)
	log.Println("code was here", singleFood.ID, err)
	if err != nil {
		log.Println("find one delivery boy error", err, singleFood.ID)
		return
	}
	if singleFood.CookedAndReady && singleBoy.CurrentLocation == utils.AtRestaurant {
		WaitToPickFromRestaurant(db, singleBoy)
	}
	if singleFood.CookedAndReady && singleBoy.CurrentLocation == utils.InTransit {
		TowardsCustomer(db, singleBoy)
	}
	if singleBoy.CurrentLocation == utils.AtCustomer {
		err = WaitToGiveCustomer(db, singleFood)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
