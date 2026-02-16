package repo

import (
	"context"
	"time"

	"concurrent-wallet-order-system/internal/config"
	"concurrent-wallet-order-system/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletRepository struct{}

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{}
}

// InsertTransaction inserts a deposit or withdraw record
func (r *WalletRepository) InsertTransaction(tx *models.WalletTransaction) error {
	collection := config.DB.Collection("wallets")

	tx.CreatedAt = time.Now()

	_, err := collection.InsertOne(context.Background(), tx)
	return err
}

// GetTransactionsByUser fetches wallet history
func (r *WalletRepository) GetTransactionsByUser(userID primitive.ObjectID) ([]models.WalletTransaction, error) {
	collection := config.DB.Collection("wallets")

	cursor, err := collection.Find(
		context.Background(),
		bson.M{"userId": userID},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var transactions []models.WalletTransaction

	for cursor.Next(context.Background()) {
		var tx models.WalletTransaction
		if err := cursor.Decode(&tx); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}
