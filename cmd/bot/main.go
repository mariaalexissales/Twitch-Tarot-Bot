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

	cardRepo, err := repository.NewCardRepository()
	if err != nil {
		log.Fatalf("Failed to initialize card repository: %v", err)
	}
	log.Println("Card database successfully loaded into memory.")

	tarotService := service.NewTarotService(cardRepo, domain.MathRandomizer{})
	bot := twitch.NewBot(cfg.TwitchBotUsername, cfg.TwitchOAuthToken, tarotService)

	log.Printf("Connecting bot %s to Twitch chat...", cfg.TwitchBotUsername)
	errCh := make(chan error, 1)
	go func() {
		errCh <- bot.Start(cfg.TwitchChannel)
	}()

	// Run until the connection fails or the host asks us to stop
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
