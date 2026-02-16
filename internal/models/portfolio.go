package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Portfolio struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"userId" json:"userId"`
	Symbol string             `bson:"symbol" json:"symbol"`
	Qty    int                `bson:"quantity" json:"quantity"`
}
