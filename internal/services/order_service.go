package services

import (
	"concurrent-wallet-order-system/internal/models"
	"concurrent-wallet-order-system/internal/repo"
	"errors"
	"fmt"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService struct {
	orderRepo      *repo.OrderRepository
	portfolioRepo  *repo.PortfolioRepository
	walletService  *WalletService
	stockService   *StockService
	mu             sync.Mutex
}

func NewOrderService(
	orderRepo *repo.OrderRepository,
	portfolioRepo *repo.PortfolioRepository,
	walletService *WalletService,
	stockService *StockService,
) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		portfolioRepo: portfolioRepo,
		walletService: walletService,
		stockService:  stockService,
	}
}


func (s *OrderService) Buy(userID primitive.ObjectID, symbol string, quantity int) (*models.Order, error) {

	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	symbol = strings.ToUpper(symbol)

	// Lock to prevent race conditions during buy
	s.mu.Lock()
	defer s.mu.Unlock()

	//  Check stock exists
	stock, err := s.stockService.GetStockBySymbol(symbol)
	if err != nil {
		return nil, errors.New("stock not found")
	}

	totalCost := float64(quantity) * stock.Price

	//  Deduct wallet balance
	err = s.walletService.Withdraw(userID, totalCost)
	if err != nil {
		return nil, err
	}

	//  Update portfolio
	fmt.Println("Calling UpsertPortfolio")

	err = s.portfolioRepo.UpsertPortfolio(userID, symbol, quantity)
	if err != nil {
		return nil, err
	}

	//  Insert order
	order := &models.Order{
		UserID:   userID,
		Symbol:   symbol,
		Type:     "BUY",
		Quantity: quantity,
		Price:    stock.Price,
	}

	err = s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) Sell(userID primitive.ObjectID, symbol string, quantity int) (*models.Order, error) {

	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	symbol = strings.ToUpper(symbol)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check stock exists
	stock, err := s.stockService.GetStockBySymbol(symbol)
	if err != nil {
		return nil, errors.New("stock not found")
	}

	//  Check portfolio
	portfolio, err := s.portfolioRepo.GetPortfolio(userID, symbol)
	if err != nil {
		return nil, errors.New("stock not owned")
	}

	if portfolio.Qty < quantity {
		return nil, errors.New("insufficient stock quantity")
	}

	totalAmount := float64(quantity) * stock.Price

	//  Add money to wallet
	err = s.walletService.Deposit(userID, totalAmount)
	if err != nil {
		return nil, err
	}

	//  Reduce portfolio quantity
	err = s.portfolioRepo.UpsertPortfolio(userID, symbol, -quantity)
	if err != nil {
		return nil, err
	}

	// Insert order
	order := &models.Order{
		UserID:   userID,
		Symbol:   symbol,
		Type:     "SELL",
		Quantity: quantity,
		Price:    stock.Price,
	}

	err = s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

