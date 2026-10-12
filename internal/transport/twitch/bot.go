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

// Commands that always get the same reply
var staticCommands = map[string]string{
	"!server":  "PROJECT//DAEMON (PvE) Project Zomboid server - IP: 198.27.75.118 | Port: 27145",
	"!discord": "Join the Discord: https://discord.gg/heXQBtYUZb",
}

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
	response, ok := b.reply(context.Background(), message.Message, message.User.DisplayName)
	if !ok {
		return
	}

	b.client.Say(message.Channel, response)
}

// reply builds the response to a chat message. The command is the first word,
// so "!tarot please" works and "!tarots" does not.
func (b *Bot) reply(ctx context.Context, text, displayName string) (string, bool) {
	fields := strings.Fields(strings.ToLower(text))
	if len(fields) == 0 {
		return "", false
	}

	switch command := fields[0]; command {
	case "!tarot":
		draw, err := b.tarotService.DrawSingle(ctx)
		if err != nil {
			log.Printf("Failed to draw card: %v", err)
			return "", false
		}

		orientation := "Upright"
		if draw.IsReversed {
			orientation = "Reversed"
		}

		return fmt.Sprintf(
			"@%s 🔮 You drew: %s (%s) - %s",
			displayName,
			draw.Card.Name,
			orientation,
			draw.Meaning,
		), true

	case "!spread":
		spread, err := b.tarotService.DrawThreeCardSpread(ctx)
		if err != nil {
			log.Printf("Failed to draw spread: %v", err)
			return "", false
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

		return fmt.Sprintf(
			"@%s 🔮 Your Spread - %s",
			displayName,
			strings.Join(parts, " | "),
		), true

	default:
		response, ok := staticCommands[command]
		return response, ok
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
