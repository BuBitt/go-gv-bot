package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// PostgresConfig struct com dados de conexão para o PostgreSQL
type PostgresConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadPostgresConfig carrega os dados do arquivo .env e os estrutura no struct Config
func LoadPostgresConfig() (*PostgresConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	return &PostgresConfig{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}, nil
}

// DiscordgoConfig struct para receber as credenciais do bot
type DiscordgoConfig struct {
	DiscordBotToken string
	GuildID         int
}

func LoadDiscordgoConfig() (*DiscordgoConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	gID, err := strconv.Atoi(os.Getenv("DISCORD_GUILD_ID"))
	if err != nil {
		return nil, fmt.Errorf("cant convert GuildID to int")
	}

	return &DiscordgoConfig{
		DiscordBotToken: os.Getenv("DISCORD_API_TOKEN"),
		GuildID:         gID,
	}, nil
}
