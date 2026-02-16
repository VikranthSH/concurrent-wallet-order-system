package handlers

import (
	"net/http"

	"concurrent-wallet-order-system/internal/services"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	stockService *services.StockService
}

func NewStockHandler(stockService *services.StockService) *StockHandler {
	return &StockHandler{
		stockService: stockService,
	}
}

type CreateStockRequest struct {
	Symbol string  `json:"symbol" binding:"required"`
	Name   string  `json:"name" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}

func (h *StockHandler) CreateStock(c *gin.Context) {
	var req CreateStockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stock, err := h.stockService.CreateStock(req.Symbol, req.Name, req.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stock)
}

func (h *StockHandler) GetAllStocks(c *gin.Context) {
	stocks, err := h.stockService.GetAllStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (h *StockHandler) GetStock(c *gin.Context) {
	symbol := c.Param("symbol")

	stock, err := h.stockService.GetStockBySymbol(symbol)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "stock not found"})
		return
	}

	c.JSON(http.StatusOK, stock)
}
