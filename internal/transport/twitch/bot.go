package twitch

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"twitch-tarot-bot/internal/service"

	"github.com/gempir/go-twitch-irc/v4"
)

type Bot struct {
	client       *twitch.Client
	tarotService *service.TarotService
}

func NewBot(username, oauthToken string, tarotService *service.TarotService) *Bot {
	client := twitch.NewClient(username, oauthToken)

	bot := &Bot{
		client:       client,
		tarotService: tarotService,
	}

	// Register event handlers
	client.OnConnect(func() {
		log.Println("Connected to Twitch chat.")
	})
	client.OnPrivateMessage(bot.handleMessage)

	return bot
}

func (b *Bot) handleMessage(message twitch.PrivateMessage) {
	ctx := context.Background()
	msgText := strings.ToLower(message.Message)

	if strings.HasPrefix(msgText, "!tarot") {
		draw, err := b.tarotService.DrawSingle(ctx)
		if err != nil {
			log.Printf("Failed to draw card: %v", err)
			return
		}

		orientation := "Upright"
		if draw.IsReversed {
			orientation = "Reversed"
		}

		response := fmt.Sprintf(
			"@%s 🔮 You drew: %s (%s) - %s",
			message.User.DisplayName,
			draw.Card.Name,
			orientation,
			draw.Meaning,
		)

		b.client.Say(message.Channel, response)
	}

	if strings.HasPrefix(msgText, "!spread") {
		spread, err := b.tarotService.DrawThreeCardSpread(ctx)
		if err != nil {
			log.Printf("Failed to draw spread: %v", err)
			return
		}

		var parts []string
		for _, draw := range spread {
			orientation := "Upright"
			if draw.IsReversed {
				orientation = "Reversed"
			}

			positionLabel := ""
			if draw.Position != nil {
				positionLabel = *draw.Position + ": "
			}
			parts = append(parts, fmt.Sprintf("%s%s (%s)", positionLabel, draw.Card.Name, orientation))
		}

		response := fmt.Sprintf(
			"@%s 🔮 Your Spread - %s",
			message.User.DisplayName,
			strings.Join(parts, " | "),
		)

		b.client.Say(message.Channel, response)
	}
}

// Start blocks until the connection is closed. It reconnects on its own after
// network drops, so a return always means the bot is done.
func (b *Bot) Start(channel string) error {
	b.client.Join(channel)

	err := b.client.Connect()
	if errors.Is(err, twitch.ErrLoginAuthenticationFailed) {
		return fmt.Errorf("twitch rejected the login, the token is expired or invalid; generate a new one and update TWITCH_OAUTH_TOKEN: %w", err)
	}
	return err
}

func (b *Bot) Stop() {
	if err := b.client.Disconnect(); err != nil {
		log.Printf("Disconnect: %v", err)
	}
}
