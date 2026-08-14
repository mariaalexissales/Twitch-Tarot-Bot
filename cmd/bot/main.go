package main

import (
	"log"
	"os"

	"twitch-tarot-bot/internal/repository"
	"twitch-tarot-bot/internal/service"
	"twitch-tarot-bot/internal/transport/twitch"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	botName := os.Getenv("TWITCH_BOT_USERNAME")
	oauthToken := os.Getenv("TWITCH_OAUTH_TOKEN")
	targetChannel := os.Getenv("TWITCH_CHANNEL")

	// 1. Load Cards into Memory
	cardRepo, err := repository.NewCardRepository("assets/data/cards.json")
	if err != nil {
		log.Fatalf("Failed to initialize card repository: %v", err)
	}
	log.Println("Card database successfully loaded into memory.")

	// 2. Initialize Service Layer
	tarotService := service.NewTarot(cardRepo)

	// 3. Initialize Twitch Bot
	bot := twitch.NewBot(botName, oauthToken, tarotService)

	// 4. Connect and Join Channel
	log.Printf("Connecting bot %s to Twitch chat...", botName)
	if err := bot.Start(targetChannel); err != nil {
		log.Fatalf("Bot error: %v", err)
	}
}
