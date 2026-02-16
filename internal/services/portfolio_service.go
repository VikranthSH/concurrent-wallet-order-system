package services

import (
	"concurrent-wallet-order-system/internal/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PortfolioService struct {
	portfolioRepo *repo.PortfolioRepository
	stockService  *StockService
}

func NewPortfolioService(
	portfolioRepo *repo.PortfolioRepository,
	stockService *StockService,
) *PortfolioService {
	return &PortfolioService{
		portfolioRepo: portfolioRepo,
		stockService:  stockService,
	}
}

type HoldingResponse struct {
	Symbol       string  `json:"symbol"`
	StockName    string  `json:"stockName"`
	Quantity     int     `json:"quantity"`
	CurrentPrice float64 `json:"currentPrice"`
	TotalValue   float64 `json:"totalValue"`
}

type PortfolioResponse struct {
	UserID              primitive.ObjectID `json:"userId"`
	Holdings            []HoldingResponse  `json:"holdings"`
	TotalPortfolioValue float64            `json:"totalPortfolioValue"`
}

func (s *PortfolioService) GetPortfolio(userID primitive.ObjectID) (*PortfolioResponse, error) {

	holdings, err := s.portfolioRepo.GetUserPortfolio(userID)
	if err != nil {
		return nil, err
	}

	var response PortfolioResponse
	response.UserID = userID

	totalValue := 0.0

	for _, h := range holdings {

		stock, err := s.stockService.GetStockBySymbol(h.Symbol)
		if err != nil {
			continue
		}

		value := float64(h.Qty) * stock.Price

		response.Holdings = append(response.Holdings, HoldingResponse{
			Symbol:       h.Symbol,
			StockName:    stock.Name,
			Quantity:     h.Qty,
			CurrentPrice: stock.Price,
			TotalValue:   value,
		})

		totalValue += value
	}

	response.TotalPortfolioValue = totalValue

	return &response, nil
}
