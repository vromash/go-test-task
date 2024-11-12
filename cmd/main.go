package main

import (
	"fmt"
	"go_test_task/internal/config"
	"go_test_task/internal/handlers"
	"go_test_task/internal/repositories"
	"go_test_task/internal/services"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.ReadConfig()
	db, err := getDatabaseClient(cfg.DB)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	balanceRepository := repositories.NewBalanceRepository(db)
	userRepository := repositories.NewUserRepository(db)

	balanceService := services.NewBalanceService(logger, balanceRepository, userRepository)

	balanceHandler := handlers.NewBalanceHandler(balanceService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/user/{userId}", func(r chi.Router) {
		r.Get("/balance", balanceHandler.GetBalance)
		r.Post("/transaction", balanceHandler.UpdateBalance)
	})

	logger.Info(fmt.Sprintf("starting http server on %d port", cfg.Port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("failed to start http server", zap.Error(err))
	}
}

func getDatabaseClient(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection with database: %v", err)
	}
	return db, nil
}
