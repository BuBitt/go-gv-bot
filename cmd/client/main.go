package main

import (
	"fmt"

	"github.com/BuBitt/gv_bot_go/cmd/client/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type User struct {
	ID       int    `db:"id"`
	Username string `db:"requester_id"`
	Email    string `db:"requester_name"`
}

func main() {
	logger.Logger.Info("Bot Launch")
	db, err := LoadPostgres()
	if err != nil {
		logger.Logger.Error("GuildID cast to int has failed", zap.Error(err))
	}

	// Example: Query data from a table
	var users []User
	err = db.Select(&users, "SELECT id, requester_id, requester_name FROM transactions")
	if err != nil {
		logger.Logger.Error("Query has failed", zap.Error(err))
	}
	// defer db.Close()

	aux := 0
	// Print user information
	for _, u := range users {
		aux += 1
		formattedString := fmt.Sprintf("ID: %d, Username: %s, Email: %s", u.ID, u.Username, u.Email)
		logger.Logger.Warn(formattedString)
		if aux == 10 {
			logger.Logger.Info("DONE BABY")
			break
		}
	}
}
