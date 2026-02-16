package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletTransaction struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"userId" json:"userId"`
	Method string `bson:"method" json:"method"`
	Amount float64 `bson:"amount" json:"amount"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}