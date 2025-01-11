package deliveryperson

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"food-delivery/models"
	"food-delivery/utils"
)

type DeliveryPersonServicesInterface interface {
	CreateDeliveryPerson(ctx context.Context, deliveryPerson models.DeliveryPerson) (*models.DeliveryPerson, error)
}

type DeliveryPersonService struct {
	Database *mongo.Database
}

func NewDeliveryPersonServicesInterface(db *mongo.Database) DeliveryPersonServicesInterface {
	return &DeliveryPersonService{Database: db}
}

func (s *DeliveryPersonService) CreateDeliveryPerson(ctx context.Context, deliveryPerson models.DeliveryPerson) (*models.DeliveryPerson, error) {
	db, err := utils.DatabaseConnect(utils.FoodDelivery)
	if err != nil {
		log.Println("database connection error", err)
		return nil, err
	}

	dfilter := primitive.M{}

	count, err := db.Collection(utils.DeliveryPersons).CountDocuments(context.Background(), dfilter)
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
			ID:              primitive.NewObjectID().Hex(),
			Name:            "boy" + fmt.Sprint(i),
			BusyStatus:      false,
			CurrentLocation: utils.Dock,
		}
		_, err = db.Collection(utils.DeliveryPersons).InsertOne(ctx, deliveryPerson)
		if err != nil {
			log.Println("DeliveryPerson InsertOne error occurred", err)
			return nil, err
		}
	}

	return &models.DeliveryPerson{}, nil
}
