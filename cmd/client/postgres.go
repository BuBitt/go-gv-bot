package main

import (
	"fmt"

	"github.com/BuBitt/gv_bot_go/cmd/client/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func LoadPostgres() (*sqlx.DB, error) {
	logger.Logger.Info("Starting PostegresSQL database connection")

	// Load PostgreSQL configuration
	config, err := LoadPostgresConfig()
	if err != nil {
		logger.Logger.Fatal("LoadPostgresConfig has failed", zap.Error(err))
	}

	// Construct the PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	// Connect to the PostgreSQL database using sqlx
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		logger.Logger.Fatal("Connection to the PostgreSQL database using sqlx has failed", zap.Error(err))
	}

	// Check if the connection to the database is successful
	err = db.Ping()
	if err != nil {
		logger.Logger.Fatal("Connection check with the database has failed", zap.Error(err))
	}

	logger.Logger.Info("Connected to the PostgreSQL database!")
	return db, nil
}
