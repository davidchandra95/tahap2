package main

import (
	"net/http"
	"tahap2/internal/config"
	"tahap2/internal/handlers"
	"tahap2/internal/middlewares"
	"tahap2/internal/repositories"
	"tahap2/internal/services"
	"tahap2/internal/workers"

	"github.com/labstack/echo/v4"
)

func main() {
	db := config.InitDB()
	e := echo.New()
	eventBus := workers.NewEventBus()

	userRepo := repositories.NewUserRepository(db)
	transRepo := repositories.NewTransactionRepo(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	transService := services.NewTransactionService(userRepo, transRepo, eventBus)
	transHandler := handlers.NewTransactionHandler(transService)

	transferWorkers := workers.NewTransactionWorker(eventBus, userRepo, transRepo)
	go transferWorkers.StartWorker()

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "pong"})
	})
	apiV1 := e.Group("/api/v1")
	apiV1.POST("/register", authHandler.Register)
	apiV1.POST("/login", authHandler.Login)
	apiV1.PUT("/profile", authHandler.UpdateProfile, middlewares.AuthMiddleware)

	apiV1.POST("/topup", transHandler.TopupHandler, middlewares.AuthMiddleware)
	apiV1.POST("/pay", transHandler.PaymentHandler, middlewares.AuthMiddleware)
	apiV1.POST("/transfer", transHandler.TransferHandler, middlewares.AuthMiddleware)
	apiV1.GET("/transactions", transHandler.GetAllTransactions, middlewares.AuthMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
