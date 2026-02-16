package services

import (
	"errors"
	"strings"

	"concurrent-wallet-order-system/internal/models"
	"concurrent-wallet-order-system/internal/repo"
)

type StockService struct {
	stockRepo *repo.StockRepository
}

func NewStockService(stockRepo *repo.StockRepository) *StockService {
	return &StockService{
		stockRepo: stockRepo,
	}
}

// Create stock
func (s *StockService) CreateStock(symbol, name string, price float64) (*models.Stock, error) {

	if price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}

	symbol = strings.ToUpper(symbol)

	// Check if stock already exists
	existing, _ := s.stockRepo.GetStockBySymbol(symbol)
	if existing != nil {
		return nil, errors.New("stock already exists")
	}

	stock := &models.Stock{
		Symbol: symbol,
		Name:   name,
		Price:  price,
	}

	err := s.stockRepo.CreateStock(stock)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *StockService) GetAllStocks() ([]models.Stock, error) {
	return s.stockRepo.GetAllStocks()
}

func (s *StockService) GetStockBySymbol(symbol string) (*models.Stock, error) {
	return s.stockRepo.GetStockBySymbol(strings.ToUpper(symbol))
}
