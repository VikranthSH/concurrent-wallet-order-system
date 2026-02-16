package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password" json:"-"`
	WalletBalance float64 `bson:"walletbalance" json:"walletbalance"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}