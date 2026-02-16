package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct{
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Symbol string `bson:"symbol" json:"symbol"`
	Name string `bson:"name" json:"name"`
	Price float64 `bson:"price" json:"price"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}