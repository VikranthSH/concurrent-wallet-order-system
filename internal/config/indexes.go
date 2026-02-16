package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes() {

	// ======================
	// Users Collection Indexes
	// ======================
	users := DB.Collection("users")

	_, err := users.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"email": 1},
		Options: options.Index().
			SetUnique(true).
			SetBackground(true),
	})
	if err != nil {
		log.Println("Failed to create users index:", err)
	}

	// ======================
	// Stocks Collection Indexes
	// ======================
	stocks := DB.Collection("stocks")

	_, err = stocks.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"symbol": 1},
		Options: options.Index().
			SetUnique(true).
			SetBackground(true),
	})
	if err != nil {
		log.Println("Failed to create stocks index:", err)
	}

	// ======================
	// Portfolio Collection Indexes
	// ======================
	portfolio := DB.Collection("portfolio")

	_, err = portfolio.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"userId": 1,
			"symbol": 1,
		},
		Options: options.Index().
			SetUnique(true).
			SetBackground(true),
	})
	if err != nil {
		log.Println("Failed to create portfolio index:", err)
	}

	// ======================
	// Orders Collection Index
	// ======================
	orders := DB.Collection("orders")

	_, err = orders.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"userId": 1},
		Options: options.Index().
			SetBackground(true),
	})
	if err != nil {
		log.Println("Failed to create orders index:", err)
	}

	log.Println("Indexes created successfully")
}
