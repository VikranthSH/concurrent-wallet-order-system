package handlers

import (
	"net/http"

	"concurrent-wallet-order-system/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PortfolioHandler struct {
	portfolioService *services.PortfolioService
}

func NewPortfolioHandler(portfolioService *services.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioService: portfolioService,
	}
}

func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	userIDParam := c.Param("userId")

	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	portfolio, err := h.portfolioService.GetPortfolio(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}
