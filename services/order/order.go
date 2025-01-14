package order

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"food-delivery/models"
	"food-delivery/utils"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, foodOrder models.FoodOrder) (*models.FoodOrder, error)
	UpdateOrder(ctx context.Context, foodOrder models.FoodOrder) (models.FoodOrder, error)
	GetOrderWithFilter(ctx context.Context, filter bson.M) ([]models.FoodOrder, error)
}

type OrderService struct {
	Database *mongo.Database
}

func NewOrderServiceInterface(db *mongo.Database) OrderServiceInterface {
	return &OrderService{Database: db}
}

func (s *OrderService) CreateOrder(ctx context.Context, foodOrder models.FoodOrder) (*models.FoodOrder, error) {
	foodOrder.ID = primitive.NewObjectID()
	foodOrder.PlacedTime = time.Now()

	_, err := s.Database.Collection(utils.Orders).InsertOne(ctx, foodOrder)
	if err != nil {
		log.Println("orders InsertOne error occurred", err)
		return nil, err
	}

	return &foodOrder, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, foodOrder models.FoodOrder) (models.FoodOrder, error) {
	return models.FoodOrder{}, nil
}

func (s *OrderService) GetOrderWithFilter(ctx context.Context, filter bson.M) ([]models.FoodOrder, error) {
	log.Println(filter)
	resultset, err := s.Database.Collection(utils.Orders).Find(ctx, bson.M{"cookassigned": false})
	if err != nil {
		log.Println("CheckOrdersRestaurant Orders find all ", err)
		return nil, err
	}
	defer resultset.Close(ctx)

	ordersList := make([]models.FoodOrder, 0)
	for resultset.Next(ctx) {
		var singleOrder models.FoodOrder
		if err := resultset.Decode(&singleOrder); err != nil {
			log.Println("CheckOrdersRestaurant Orders decode error ", err)
			return []models.FoodOrder{}, nil
		}
		ordersList = append(ordersList, singleOrder)
	}

	return ordersList, nil
}
