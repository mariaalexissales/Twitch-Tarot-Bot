package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"twitch-tarot-bot/internal/config"
	"twitch-tarot-bot/internal/domain"
	"twitch-tarot-bot/internal/repository"
	"twitch-tarot-bot/internal/service"
	"twitch-tarot-bot/internal/transport/twitch"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// 1. Load Cards into Memory
	cardRepo, err := repository.NewCardRepository()
	if err != nil {
		log.Fatalf("Failed to initialize card repository: %v", err)
	}
	log.Println("Card database successfully loaded into memory.")

	// 2. Initialize Service Layer
	tarotService := service.NewTarotService(cardRepo, domain.CryptoRandomizer{})

	// 3. Initialize Twitch Bot
	bot := twitch.NewBot(cfg.TwitchBotUsername, cfg.TwitchOAuthToken, tarotService)

	// 4. Connect and Join Channel
	log.Printf("Connecting bot %s to Twitch chat...", cfg.TwitchBotUsername)
	errCh := make(chan error, 1)
	go func() {
		errCh <- bot.Start(cfg.TwitchChannel)
	}()

	// 5. Run until the connection fails or the host asks us to stop
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		log.Fatalf("Bot error: %v", err)
	case <-ctx.Done():
		log.Println("Shutting down...")
		bot.Stop()
	}
}
