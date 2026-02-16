package services

import (
	"errors"
	"time"

	"concurrent-wallet-order-system/internal/models"
	"concurrent-wallet-order-system/internal/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repo.UserRepository
}

func NewUserService(userRepo *repo.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Register(name, email, password string) (*models.User, error) {

	// Check if user already exists
	existingUser, _ := s.userRepo.GetUserByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:          name,
		Email:         email,
		Password:      string(hashedPassword),
		WalletBalance: 0,
		CreatedAt:     time.Now(),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(email, password string) (*models.User, error) {

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *UserService) GetUserByID(userID primitive.ObjectID) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}
