package services

import (
	"errors"
	"sync"

	"concurrent-wallet-order-system/internal/models"
	"concurrent-wallet-order-system/internal/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletService struct {
	userRepo   *repo.UserRepository
	walletRepo *repo.WalletRepository
	mu         sync.Mutex
}

func NewWalletService(
	userRepo *repo.UserRepository,
	walletRepo *repo.WalletRepository,
) *WalletService {
	return &WalletService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

func (s *WalletService) Deposit(userID primitive.ObjectID, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	newBalance := user.WalletBalance + amount

	err = s.userRepo.UpdateWalletBalance(userID, newBalance)
	if err != nil {
		return err
	}

	tx := &models.WalletTransaction{
		UserID: userID,
		Method: "deposit",
		Amount: amount,
	}

	return s.walletRepo.InsertTransaction(tx)
}
func (s *WalletService) Withdraw(userID primitive.ObjectID, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.WalletBalance < amount {
		return errors.New("insufficient balance")
	}

	newBalance := user.WalletBalance - amount

	err = s.userRepo.UpdateWalletBalance(userID, newBalance)
	if err != nil {
		return err
	}

	tx := &models.WalletTransaction{
		UserID: userID,
		Method: "withdraw",
		Amount: amount,
	}

	return s.walletRepo.InsertTransaction(tx)
}

func (s *WalletService) GetBalance(userID primitive.ObjectID) (float64, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return 0, err
	}

	return user.WalletBalance, nil
}

func (s *WalletService) GetHistory(userID primitive.ObjectID) ([]models.WalletTransaction, error) {
	return s.walletRepo.GetTransactionsByUser(userID)
}
