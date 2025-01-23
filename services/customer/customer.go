package customer

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"food-delivery/models"
	"food-delivery/services/order"
	"food-delivery/utils"
)

type CustomerServiceInterface interface {
	CreateCustomer(ctx context.Context, customer models.Customer) (*models.Customer, error)
	GetCustomerByID(ctx context.Context, customerID string) (models.Customer, error)
	GetCustomerByPhone(ctx context.Context, phone string) (*models.Customer, error)
	UpdateCustomerOrders(ctx context.Context, customer models.Customer, foodOrder models.FoodOrder) (*models.Customer, error)
	UpdateCustomerOrderCustom(ctx context.Context, findFilter, updateSet bson.M) error
	GetAllCustomers(ctx context.Context) ([]models.Customer, error)
	GetAllWaitingCustomersCustom(ctx context.Context, findFilter bson.M) ([]models.Customer, error)
	GetCustomerOrders(ctx context.Context, customerID primitive.ObjectID) ([]models.FoodOrder, error)
}

type CustomerService struct {
	Database     *mongo.Database
	OrderService order.OrderServiceInterface
}

func NewCustomerServiceInterface(db *mongo.Database) CustomerServiceInterface {
	return &CustomerService{Database: db, OrderService: order.NewOrderServiceInterface(db)}
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

func (s *CustomerService) UpdateCustomerOrders(ctx context.Context, _ models.Customer, foodOrder models.FoodOrder) (*models.Customer, error) {
	customerObjID, err := primitive.ObjectIDFromHex(foodOrder.CustomerID)
	if err != nil {
		log.Println("invalid customer ID", err)
		return nil, err
	}

	singleCustomerOrder := models.CustomerOrders{
		OrderName:      foodOrder.FoodName,
		CurrentOrderID: foodOrder.ID.Hex(),
		PlacedTime:     foodOrder.PlacedTime,
	}

	filter := bson.M{"id": customerObjID}
	update := bson.M{"$push": bson.M{"orders": singleCustomerOrder}}
	opts := options.Update().SetUpsert(true)

	_, err = s.Database.Collection(utils.Customers).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		fmt.Println("Error updating document with latest order:", err)
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

func (s *CustomerService) GetAllCustomers(ctx context.Context) ([]models.Customer, error) {
	cursor, err := s.Database.Collection(utils.Customers).Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("error getting customer list:", err)
		return []models.Customer{}, err
	}

	customers := make([]models.Customer, 0)
	for cursor.Next(ctx) {
		singleCustomer := models.Customer{}
		if err = cursor.Decode(&singleCustomer); err != nil {
			log.Println("GetAllCustomers decode error ", err)
			return []models.Customer{}, nil
		}
		customers = append(customers, singleCustomer)
	}
	if len(customers) == 0 {
		return []models.Customer{}, mongo.ErrNoDocuments
	}

	return customers, nil
}

func (s *CustomerService) GetCustomerOrders(ctx context.Context, customerID primitive.ObjectID) ([]models.FoodOrder, error) {
	findFilter := bson.M{"customerid": customerID.Hex(), "cookedandready": true}
	orders, err := s.OrderService.GetOrderWithFilter(ctx, findFilter)
	if err != nil {
		fmt.Println("error getting customer list:", err)
		return []models.FoodOrder{}, err
	}
	return orders, nil
}

func (s *CustomerService) GetAllWaitingCustomersCustom(ctx context.Context, findFilter bson.M) ([]models.Customer, error) {
	cursor, err := s.Database.Collection(utils.Customers).Find(ctx, findFilter)
	if err != nil {
		fmt.Println("error getting customer list:", err)
		return []models.Customer{}, err
	}

	customers := make([]models.Customer, 0)
	for cursor.Next(ctx) {
		singleCustomer := models.Customer{}
		if err = cursor.Decode(&singleCustomer); err != nil {
			log.Println("GetAllCustomers decode error ", err)
			return []models.Customer{}, nil
		}
		customers = append(customers, singleCustomer)
	}

	if len(customers) == 0 {
		return []models.Customer{}, mongo.ErrNoDocuments
	}

	return customers, nil
}
