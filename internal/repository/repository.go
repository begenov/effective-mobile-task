package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/begenov/effective-mobile-task/internal/logger"
	"github.com/begenov/effective-mobile-task/internal/model"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

type UserRepository interface {
	Create(context.Context, model.User) error
	Update(context.Context, model.User) error
	Delete(context.Context, int64) error
	GetUsers(ctx context.Context, userID, limit, offset *int64, gender *int, nationality *string) ([]model.User, error)
}

func NewPostgresDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info(connString)
	err = db.Ping()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	direction := "up"

	err = ApplyMigrations(db, direction)
	if err != nil {
		fmt.Println("Error applying migrations:", err)
		return nil, err
	}

	return db, nil
}

func ApplyMigrations(db *sql.DB, direction string) error {
	migrationsDir := "migration"

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if direction == "up" && filepath.Ext(file.Name()) == ".up.sql" {
			err := applyMigration(db, filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return err
			}
		} else if direction == "down" && filepath.Ext(file.Name()) == ".down.sql" {
			err := applyMigration(db, filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func applyMigration(db *sql.DB, filePath string) error {
	migrationSQL, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		return err
	}

	fmt.Printf("Applied migration: %s\n", filePath)
	return nil
}
