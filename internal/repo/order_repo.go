package repo

import (
	"context"
	"time"

	"concurrent-wallet-order-system/internal/config"
	"concurrent-wallet-order-system/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	collection := config.DB.Collection("orders")

	order.CreatedAt = time.Now()

	result, err := collection.InsertOne(context.Background(), order)
	if err != nil {
		return err
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}
