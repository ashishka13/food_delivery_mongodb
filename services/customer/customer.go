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
	GetCustomer(ctx context.Context, customer models.Customer) (*models.Customer, error)
	UpdateCustomerOrder(ctx context.Context, customer models.Customer, foodOrder models.FoodOrder) (*models.Customer, error)
}

type CustomerService struct {
	Database *mongo.Database
}

func NewCustomerServiceInterface(db *mongo.Database) CustomerServiceInterface {
	return &CustomerService{Database: db}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer models.Customer) (*models.Customer, error) {
	customer.ID = primitive.NewObjectID().Hex()
	_, err := s.Database.Collection(utils.Customers).InsertOne(ctx, customer)
	if err != nil {
		log.Println("new customer InsertOne error occurred", err)
		return nil, err
	}

	return &customer, nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, customer models.Customer) (*models.Customer, error) {
	customerFound := models.Customer{}
	err := s.Database.Collection(utils.Customers).FindOne(context.Background(), bson.M{"id": customer.ID}).Decode(&customerFound)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("customer not found")
			return nil, err
		}
		log.Println("createOrder customer find one error", err)
		return nil, err
	}
	return &customerFound, nil
}

func (s *CustomerService) UpdateCustomerOrder(ctx context.Context, customer models.Customer, foodOrder models.FoodOrder) (*models.Customer, error) {
	filter := bson.M{"id": customer.ID}
	update := bson.M{"$set": bson.M{"currentorderid": foodOrder.ID, "orderplaced": true, "placedtime": foodOrder.PlacedTime}}

	res, err := s.Database.Collection(utils.Customers).UpdateOne(ctx, filter, update)
	if err != nil || res.ModifiedCount == 0 {
		log.Println("find single customer error", err)
		return nil, err
	}

	return &models.Customer{}, nil
}
