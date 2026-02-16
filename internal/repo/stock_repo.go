package repo

import (
	"context"
	"time"

	"concurrent-wallet-order-system/internal/config"
	"concurrent-wallet-order-system/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockRepository struct{}

func NewStockRepository() *StockRepository {
	return &StockRepository{}
}

// Create stock
func (r *StockRepository) CreateStock(stock *models.Stock) error {
	collection := config.DB.Collection("stocks")

	stock.CreatedAt = time.Now()

	result, err := collection.InsertOne(context.Background(), stock)
	if err != nil {
		return err
	}

	stock.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Get all stocks
func (r *StockRepository) GetAllStocks() ([]models.Stock, error) {
	collection := config.DB.Collection("stocks")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var stocks []models.Stock

	for cursor.Next(context.Background()) {
		var stock models.Stock
		if err := cursor.Decode(&stock); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// Get stock by symbol
func (r *StockRepository) GetStockBySymbol(symbol string) (*models.Stock, error) {
	collection := config.DB.Collection("stocks")

	var stock models.Stock
	err := collection.FindOne(
		context.Background(),
		bson.M{"symbol": symbol},
	).Decode(&stock)

	if err != nil {
		return nil, err
	}

	return &stock, nil
}
