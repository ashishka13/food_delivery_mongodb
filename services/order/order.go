package order

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"food-delivery/models"
	"food-delivery/utils"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, foodOrder models.FoodOrder) (*models.FoodOrder, error)
	UpdateOrder(ctx context.Context, foodOrder models.FoodOrder) (*models.FoodOrder, error)
}

type OrderService struct {
	Database *mongo.Database
}

func NewOrderServiceInterface(db *mongo.Database) OrderServiceInterface {
	return &OrderService{Database: db}
}

func (s *OrderService) CreateOrder(ctx context.Context, foodOrder models.FoodOrder) (*models.FoodOrder, error) {
	foodOrder.ID = primitive.NewObjectID().Hex()
	foodOrder.PlacedTime = time.Now()

	_, err := s.Database.Collection(utils.Orders).InsertOne(ctx, foodOrder)
	if err != nil {
		log.Println("orders InsertOne error occurred", err)
		return nil, err
	}

	return &foodOrder, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, foodOrder models.FoodOrder) (*models.FoodOrder, error) {
	return &models.FoodOrder{}, nil
}
