package deliveryperson

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"food-delivery/models"
	"food-delivery/utils"
)

type DeliveryPersonServicesInterface interface {
	CreateDeliveryPerson(ctx context.Context, deliveryPerson models.DeliveryPerson) (*models.DeliveryPerson, error)
	UpdateDeliveryPersonCustom(ctx context.Context, findFilter, updateSet bson.M) error
	GetDeliveryPersonCustom(ctx context.Context, findFilter bson.M) (models.DeliveryPerson, error)
}

type DeliveryPersonService struct {
	Database *mongo.Database
}

func NewDeliveryPersonServicesInterface(db *mongo.Database) DeliveryPersonServicesInterface {
	return &DeliveryPersonService{Database: db}
}

func (s *DeliveryPersonService) CreateDeliveryPerson(ctx context.Context, deliveryPerson models.DeliveryPerson) (*models.DeliveryPerson, error) {
	log.Println("code was here")
	dfilter := primitive.M{}

	count, err := s.Database.Collection(utils.DeliveryPersons).CountDocuments(context.Background(), dfilter)
	if err != nil {
		log.Println("CreateDeliveryPerson count error DeliveryPerson", err)
		return nil, err
	}
	if count > 0 {
		log.Println("delivery boys already present")
		return nil, nil
	}

	for i := 0; i < 4; i++ {
		deliveryPerson := models.DeliveryPerson{
			ID:              primitive.NewObjectID(),
			Name:            "boy" + fmt.Sprint(i),
			BusyStatus:      false,
			CurrentLocation: utils.Dock,
		}
		_, err = s.Database.Collection(utils.DeliveryPersons).InsertOne(ctx, deliveryPerson)
		if err != nil {
			log.Println("DeliveryPerson InsertOne error occurred", err)
			return nil, err
		}
	}

	return &models.DeliveryPerson{}, nil
}

func (s *DeliveryPersonService) UpdateDeliveryPersonCustom(ctx context.Context, findFilter, updateSet bson.M) error {
	updateresult, err := s.Database.Collection(utils.DeliveryPersons).UpdateOne(context.Background(), findFilter, updateSet)
	if err != nil {
		log.Println("error occurred while updating the delivery person", err)
		return err
	}
	if updateresult.ModifiedCount == 0 {
		return errors.New("no documents were updated")
	}
	return nil
}

func (s *DeliveryPersonService) GetDeliveryPersonCustom(ctx context.Context, findFilter bson.M) (models.DeliveryPerson, error) {
	deliveryPersonFound := models.DeliveryPerson{}
	err := s.Database.Collection(utils.DeliveryPersons).FindOne(ctx, findFilter).Decode(&deliveryPersonFound)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("error occurred while getting the delivery person", err)
		return models.DeliveryPerson{}, err
	}
	return deliveryPersonFound, nil
}

/*
dfilter := primitive.M{"_id": singleBoy.ID}
	dupdate := primitive.M{"$set": primitive.M{"currentlocation": utils.InTransit, "busystatus": true, "currentorderid": singleFood.ID}}
	updateresult, err := db.Collection(utils.DeliveryPersons).UpdateOne(context.Background(), dfilter, dupdate)

	dupdate = primitive.M{"$set": primitive.M{"currentlocation": utils.AtRestaurant}}
	_, err = db.Collection(utils.DeliveryPersons).UpdateOne(context.Background(), dfilter, dupdate)

	dfilter := primitive.M{"_id": singleBoy.ID}
	dupdate := primitive.M{"$set": primitive.M{"currentlocation": utils.InTransit}}
	_, err := db.Collection(utils.DeliveryPersons).UpdateOne(context.Background(), dfilter, dupdate)

	dfilter := primitive.M{"_id": singleBoy.ID}
	dupdate := primitive.M{"$set": primitive.M{"currentlocation": utils.AtCustomer}}
	_, err := db.Collection(utils.DeliveryPersons).UpdateOne(context.Background(), dfilter, dupdate)

		dfilter := primitive.M{"currentorderid": singleFood.ID}
	dupdate := primitive.M{"$set": primitive.M{"busystatus": false, "currentorderid": "", "currentlocation": utils.Dock}}
	_, err = db.Collection(utils.DeliveryPersons).UpdateOne(context.Background(), dfilter, dupdate)
*/
