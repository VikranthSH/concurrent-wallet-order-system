package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Symbol    string             `bson:"symbol" json:"symbol"`
	Type      string             `bson:"type" json:"type"` // BUY or SELL
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
