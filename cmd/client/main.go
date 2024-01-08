package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/BuBitt/gv_bot_go/cmd/client/logger"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type User struct {
	ID       int    `db:"id"`
	Username string `db:"requester_id"`
	Email    string `db:"requester_name"`
}

var (
	DiscordGoConfig, err = LoadDiscordgoConfig()
	Token                = DiscordGoConfig.DiscordBotToken
	GuildID              = DiscordGoConfig.GuildID
	waitingUsers         = make(map[string]func(*discordgo.Session, *discordgo.MessageCreate))
)

func main() {
	logger.Logger.Info("Bot Launch")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		logger.Logger.Error("Error creating Discord session", zap.Error(err))
		return
	}
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	if err != nil {
		logger.Logger.Error("GuildID cast to int has failed", zap.Error(err))
	}

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		logger.Logger.Error("error opening connection", zap.Error(err))
		return
	}

	// db, err := LoadPostgres()
	// Example: Query data from a table
	// var users []User
	// err = db.Select(&users, "SELECT id, requester_id, requester_name FROM transactions")
	// if err != nil {
	// 	logger.Logger.Error("Query has failed", zap.Error(err))
	// }
	// defer db.Close()
	//
	// aux := 0
	// // Print user information
	// for _, u := range users {
	// 	aux += 1
	// 	formattedString := fmt.Sprintf("ID: %d, Username: %s, Email: %s", u.ID, u.Username, u.Email)
	// 	logger.Logger.Warn(formattedString)
	// 	if aux == 10 {
	// 		logger.Logger.Info("DONE BABY")
	// 		break
	// 	}
	// }

	// Wait here until CTRL-C or other term signal is received.
	logger.Logger.Warn("Bot is now running. Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	logger.Logger.Warn("Shutting down...")
	dg.Close()
	logger.Logger.Warn("Shutting down complete")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if message.Author.ID == session.State.User.ID {
		return
	}
	msg := fmt.Sprintf("Author: %s | Message: %s", message.Author.Username, message.Content)
	logger.Logger.Info(msg)
	// Check if the user is waiting for a response
	if waiting, ok := waitingUsers[message.Author.ID]; ok {
		// Execute the function associated with the user
		waiting(session, message)
		// Remove the user from the waiting list
		if message.Content == "!Stop" {
			delete(waitingUsers, message.Author.ID)
		}
		return
	}

	// Handle your regular commands here
	if strings.HasPrefix(message.Content, "!command") {
		// Triggered command, now wait for another message from the user
		waitForNextMessage(session, message.Author.ID, func(session *discordgo.Session, message *discordgo.MessageCreate) {
			// Process the next message
			session.ChannelMessageSend(message.ChannelID, message.Content)
			fstrng := fmt.Sprintf("Received another message from %s: %s", message.Author.Username, message.Content)
			logger.Logger.Info(fstrng)
		})
	}
}

func waitForNextMessage(session *discordgo.Session, userID string, callback func(*discordgo.Session, *discordgo.MessageCreate)) {
	// Set the user in the waiting list with the associated callback function
	waitingUsers[userID] = callback
	// You can add a timeout mechanism if needed
}
