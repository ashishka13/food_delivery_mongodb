package customer

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"food-delivery/models"
	"food-delivery/utils"
)

type CustomerServiceInterface interface {
	CreateCustomer(ctx context.Context, customer models.Customer) (*models.Customer, error)
	GetCustomerByID(ctx context.Context, customerID string) (models.Customer, error)
	GetCustomerByPhone(ctx context.Context, phone string) (*models.Customer, error)
	UpdateCustomerOrderID(ctx context.Context, customer models.Customer, foodOrder models.FoodOrder) (*models.Customer, error)
	UpdateCustomerOrderCustom(ctx context.Context, findFilter, updateSet bson.M) error
}

type CustomerService struct {
	Database *mongo.Database
}

func NewCustomerServiceInterface(db *mongo.Database) CustomerServiceInterface {
	return &CustomerService{Database: db}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer models.Customer) (*models.Customer, error) {
	customer.ID = primitive.NewObjectID()
	_, err := s.Database.Collection(utils.Customers).InsertOne(ctx, customer)
	if err != nil {
		log.Println("new customer InsertOne error occurred", err)
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) GetCustomerByID(ctx context.Context, customerID string) (models.Customer, error) {
	customerObjectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		log.Println("invalid customerID entered")
		return models.Customer{}, err
	}

	customerFound := models.Customer{}
	err = s.Database.Collection(utils.Customers).FindOne(context.Background(), bson.M{"id": customerObjectID}).Decode(&customerFound)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("customer not found")
			return models.Customer{}, err
		}
		log.Println("createOrder customer find one error", err)
		return models.Customer{}, err
	}
	return customerFound, nil
}

func (s *CustomerService) GetCustomerByPhone(ctx context.Context, phone string) (*models.Customer, error) {
	customerFound := models.Customer{}
	err := s.Database.Collection(utils.Customers).FindOne(context.Background(), bson.M{"phone": phone}).Decode(&customerFound)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("customer not found with given phone")
			return nil, err
		}
		log.Println("createOrder customer find one error", err)
		return nil, err
	}
	return &customerFound, nil
}

func (s *CustomerService) UpdateCustomerOrderID(ctx context.Context, customer models.Customer, foodOrder models.FoodOrder) (*models.Customer, error) {
	filter := bson.M{"id": customer.ID}
	update := bson.M{"$set": bson.M{"currentorderid": foodOrder.ID, "orderplaced": true, "placedtime": foodOrder.PlacedTime}}

	res, err := s.Database.Collection(utils.Customers).UpdateOne(ctx, filter, update)
	if err != nil || res.ModifiedCount == 0 {
		log.Println("update single customer error", err)
		return nil, err
	}

	return &models.Customer{}, nil
}

func (s *CustomerService) UpdateCustomerOrderCustom(ctx context.Context, findFilter, updateSet bson.M) error {
	res, err := s.Database.Collection(utils.Customers).UpdateOne(ctx, findFilter, updateSet)
	if err != nil || res.ModifiedCount == 0 {
		log.Println("update single customer error", err)
		return err
	}

	return nil
}

/*
cfilter := bson.M{"currentorderid": singleFood.ID}
	cupdate := bson.M{"$set": bson.M{"inprocess": true}}
	_, err = db.Collection(utils.Customers).UpdateOne(context.Background(), cfilter, cupdate)

	cfilter := primitive.M{"currentorderid": singleFood.ID}
	cupdate := primitive.M{"$set": primitive.M{"deliverypersonname": singleBoy.Name}}
	updateresult, err = db.Collection(utils.Customers).UpdateOne(context.Background(), cfilter, cupdate)

	cfilter := primitive.M{"currentorderid": singleFood.ID}
	cupdate := primitive.M{"$set": primitive.M{"receivetime": time.Now(), "currentorderid": "", "orderplaced": false, "inprocess": false}}
	_, err = db.Collection(utils.Customers).UpdateOne(context.Background(), cfilter, cupdate)
*/
