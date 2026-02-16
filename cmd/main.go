package main

import (
	"log"

	"concurrent-wallet-order-system/internal/config"
	"concurrent-wallet-order-system/internal/handlers"
	"concurrent-wallet-order-system/internal/repo"
	"concurrent-wallet-order-system/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {

	// =============================
	// Connect to MongoDB
	// =============================
	mongoURI := "mongodb://localhost:27017"
	dbName := "wallet_order_system"

	config.ConnectMongo(mongoURI, dbName)
	config.CreateIndexes()


	// Repositories
	userRepo := repo.NewUserRepository()
	walletRepo := repo.NewWalletRepository()
	stockRepo := repo.NewStockRepository()
	orderRepo := repo.NewOrderRepository()
	portfolioRepo := repo.NewPortfolioRepository()

	// Services
	userService := services.NewUserService(userRepo)
	walletService := services.NewWalletService(userRepo, walletRepo)
	stockService := services.NewStockService(stockRepo)
	orderService := services.NewOrderService(
		orderRepo,
		portfolioRepo,
		walletService,
		stockService,
	)
	portfolioService := services.NewPortfolioService(portfolioRepo, stockService)


	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	walletHandler := handlers.NewWalletHandler(walletService)
	stockHandler := handlers.NewStockHandler(stockService)
	orderHandler := handlers.NewOrderHandler(orderService)
	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)


	// =============================
	// Setup Router
	// =============================
	router := gin.Default()

	// User Routes
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	// Wallet Routes
	router.POST("/wallet/deposit", walletHandler.Deposit)
	router.POST("/wallet/withdraw", walletHandler.Withdraw)
	router.GET("/wallet/balance/:userId", walletHandler.GetBalance)
	router.GET("/wallet/history/:userId", walletHandler.GetHistory)
	router.GET("/users", userHandler.GetAllUsers)
	router.GET("/users/:userId", userHandler.GetUser)
	router.POST("/stocks", stockHandler.CreateStock)
	router.GET("/stocks", stockHandler.GetAllStocks)
	router.GET("/stocks/:symbol", stockHandler.GetStock)
	router.POST("/orders/buy", orderHandler.Buy)
	router.POST("/orders/sell", orderHandler.Sell)
	router.GET("/portfolio/:userId", portfolioHandler.GetPortfolio)

	// =============================
	//  Start Server
	// =============================
	log.Println("Server running on port 8080")
	router.Run(":8080")
}
