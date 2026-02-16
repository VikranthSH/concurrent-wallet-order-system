package repo

import (
	"context"

	"concurrent-wallet-order-system/internal/config"
	"concurrent-wallet-order-system/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PortfolioRepository struct{}

func NewPortfolioRepository() *PortfolioRepository {
	return &PortfolioRepository{}
}

func (r *PortfolioRepository) GetPortfolio(userID primitive.ObjectID, symbol string) (*models.Portfolio, error) {
	collection := config.DB.Collection("portfolio")

	var p models.Portfolio
	err := collection.FindOne(
		context.Background(),
		bson.M{"userId": userID, "symbol": symbol},
	).Decode(&p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
func (r *PortfolioRepository) UpsertPortfolio(userID primitive.ObjectID, symbol string, qty int) error {
	collection := config.DB.Collection("portfolio")

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"userId": userID, "symbol": symbol},
		bson.M{
			"$inc": bson.M{"quantity": qty},
			"$setOnInsert": bson.M{
				"userId": userID,
				"symbol": symbol,
			},
		},
		options.Update().SetUpsert(true),
	)

	return err
}

func (r *PortfolioRepository) GetUserPortfolio(userID primitive.ObjectID) ([]models.Portfolio, error) {
	collection := config.DB.Collection("portfolio")

	cursor, err := collection.Find(
		context.Background(),
		bson.M{"userId": userID},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var holdings []models.Portfolio

	for cursor.Next(context.Background()) {
		var p models.Portfolio
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		holdings = append(holdings, p)
	}

	return holdings, nil
}

func (r *PortfolioRepository) GetPortfolioWithAggregation(userID primitive.ObjectID) (bson.M, error) {

	collection := config.DB.Collection("portfolio")

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"userId", userID}}}},
		{{"$lookup", bson.D{
			{"from", "stocks"},
			{"localField", "symbol"},
			{"foreignField", "symbol"},
			{"as", "stockInfo"},
		}}},
		{{"$unwind", "$stockInfo"}},
		{{"$project", bson.D{
			{"symbol", 1},
			{"quantity", 1},
			{"stockName", "$stockInfo.name"},
			{"currentPrice", "$stockInfo.price"},
			{"totalValue", bson.D{
				{"$multiply", bson.A{"$quantity", "$stockInfo.price"}},
			}},
		}}},
		{{"$group", bson.D{
			{"_id", nil},
			{"holdings", bson.D{{"$push", "$$ROOT"}}},
			{"totalPortfolioValue", bson.D{{"$sum", "$totalValue"}}},
		}}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return bson.M{
			"userId":              userID,
			"holdings":            []bson.M{},
			"totalPortfolioValue": 0,
		}, nil
	}

	return results[0], nil
}
