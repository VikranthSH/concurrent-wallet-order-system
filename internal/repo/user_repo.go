package repo

import (
	"context"
	"time"

	"concurrent-wallet-order-system/internal/config"
	"concurrent-wallet-order-system/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// CreateUser inserts a new user
func (r *UserRepository) CreateUser(user *models.User) error {
	collection := config.DB.Collection("users")

	user.CreatedAt = time.Now()
	user.WalletBalance = 0

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	// IMPORTANT: assign inserted ID back to user struct
	user.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}
// GetUserByEmail finds user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	collection := config.DB.Collection("users")

	var user models.User
	err := collection.FindOne(
		context.Background(),
		bson.M{"email": email},
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID finds user by ID
func (r *UserRepository) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	collection := config.DB.Collection("users")

	var user models.User
	err := collection.FindOne(
		context.Background(),
		bson.M{"_id": id},
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateWalletBalance updates user's wallet balance
func (r *UserRepository) UpdateWalletBalance(userID primitive.ObjectID, newBalance float64) error {
	collection := config.DB.Collection("users")

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"walletBalance": newBalance}},
	)

	return err
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	collection := config.DB.Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
