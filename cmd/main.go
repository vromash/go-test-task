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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.ReadConfig()
	db, err := getDatabaseClient(cfg.DB)
	if err != nil {
		panic(err)
	}

	balanceRepository := repositories.NewBalanceRepository(db)
	userRepository := repositories.NewUserRepository(db)

	balanceService := services.NewBalanceService(balanceRepository, userRepository)

	balanceHandler := handlers.NewBalanceHandler(balanceService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/user/{userId}", func(r chi.Router) {
		r.Get("/balance", balanceHandler.GetBalance)
		r.Post("/transaction", balanceHandler.UpdateBalance)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}

func getDatabaseClient(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection with database: %v", err)
	}
	return db, nil
}
