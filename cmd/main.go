package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/begenov/effective-mobile-task/internal/handler"
	"github.com/begenov/effective-mobile-task/internal/logger"
	"github.com/begenov/effective-mobile-task/internal/repository"
	"github.com/begenov/effective-mobile-task/internal/repository/postgre"
	"github.com/begenov/effective-mobile-task/internal/service"
	"github.com/go-chi/chi"
)

func main() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := repository.NewPostgresDB(databaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to the database: ", err)
	}

	userRepository := postgre.New(db)
	userService := service.New(userRepository)
	userHandler := handler.New(userService)

	router := chi.NewRouter()
	router.Mount("/api", userHandler.Router())

	server := &http.Server{
		Addr:    os.Getenv("APP_PORT"),
		Handler: router,
	}

	go func() {
		logger.Info("Server is running on port ", os.Getenv("APP_PORT"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start the server: ", err)
		}
	}()

	select {
	case <-signalCh:
		logger.Info("Received interrupt signal. Shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Error during server shutdown: ", err)
		}

		logger.Info("Server gracefully stopped")
	}
}
