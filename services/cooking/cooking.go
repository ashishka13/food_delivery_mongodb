package cooking

import (
	"context"
	"errors"
	"food-delivery/utils"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CookingServiceInterface interface {
	UpdateFoodOrderCustom(ctx context.Context, findFilter, updateSet bson.M) error
}

type CookingService struct {
	Database *mongo.Database
}

func NewCookingServiceInterface(db *mongo.Database) CookingServiceInterface {
	return &CookingService{Database: db}
}

func (s *CookingService) UpdateFoodOrderCustom(ctx context.Context, findFilter, updateSet bson.M) error {
	updateResult, err := s.Database.Collection(utils.Orders).UpdateOne(ctx, findFilter, updateSet)
	if err != nil {
		log.Println("error occurred while updating the food order", err)
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("no documents were updated")
	}
	return nil
}

/*
ofilter := bson.M{"id": singleFood.ID}
	oupdate := bson.M{"$set": bson.M{"cookassigned": true}}
	_, err := db.Collection(utils.Orders).UpdateOne(context.Background(), ofilter, oupdate)

	oupdate = bson.M{"$set": bson.M{"cookedandready": true}}
	_, err = db.Collection(utils.Orders).UpdateOne(context.Background(), ofilter, oupdate)


	ofilter := primitive.M{"id": singleFood.ID}
	oupdate := primitive.M{"$set": primitive.M{"deliverypersonassigned": true}}
	_, err = db.Collection(utils.Orders).UpdateOne(context.Background(), ofilter, oupdate)

	ofilter := primitive.M{"id": singleFood.ID}
	oupdate := primitive.M{"$set": primitive.M{"completestatus": true}}
	_, err = db.Collection(utils.Orders).UpdateOne(context.Background(), ofilter, oupdate)
*/
