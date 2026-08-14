package twitch

import (
	"fmt"
	"strings"

	"twitch-tarot-bot/internal/service"

	"github.com/gempir/go-twitch-irc/v4"
)

type Bot struct {
	client       *twitch.Client
	tarotService *service.Tarot
}

func NewBot(username, oauthToken string, tarotService *service.Tarot) *Bot {
	client := twitch.NewClient(username, oauthToken)

	bot := &Bot{
		client:       client,
		tarotService: tarotService,
	}

	// Register event handlers
	client.OnPrivateMessage(bot.handleMessage)

	return bot
}

func (b *Bot) handleMessage(message twitch.PrivateMessage) {
	msgText := strings.ToLower(message.Message)

	if strings.HasPrefix(msgText, "!tarot") {
		draw := b.tarotService.DrawSingle()

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
		spread := b.tarotService.DrawThreeCardSpread()

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

func (b *Bot) Start(channel string) error {
	b.client.Join(channel)
	return b.client.Connect()
}
