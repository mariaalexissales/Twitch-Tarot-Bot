package twitch

import (
	"context"
	"strings"
	"testing"

	"twitch-tarot-bot/internal/repository"
	"twitch-tarot-bot/internal/service"
)

// MockRandomizer always returns 0 for deterministic testing
type MockRandomizer struct{}

func (m MockRandomizer) IntN(n int) int { return 0 }

func newTestBot(t *testing.T) (*Bot, *repository.CardRepository) {
	t.Helper()
	repo, err := repository.NewCardRepository()
	if err != nil {
		t.Fatalf("failed to load cards: %v", err)
	}
	return &Bot{tarotService: service.NewTarotService(repo, MockRandomizer{})}, repo
}

func TestReplyStaticCommands(t *testing.T) {
	bot, _ := newTestBot(t)

	tests := []struct {
		message string
		want    []string
	}{
		{"!server", []string{"PROJECT//DAEMON", "198.27.75.118", "27145"}},
		{"!SERVER", []string{"198.27.75.118", "27145"}},
		{"!server now", []string{"198.27.75.118", "27145"}},
		{"  !server", []string{"198.27.75.118", "27145"}},
		{"!discord", []string{"discord.gg/heXQBtYUZb"}},
		{"!Discord link?", []string{"discord.gg/heXQBtYUZb"}},
	}

	for _, tt := range tests {
		response, ok := bot.reply(context.Background(), tt.message, "Viewer")
		if !ok {
			t.Errorf("%q: expected a reply", tt.message)
			continue
		}
		for _, want := range tt.want {
			if !strings.Contains(response, want) {
				t.Errorf("%q: expected %q in reply %q", tt.message, want, response)
			}
		}
	}
}

func TestReplyTarot(t *testing.T) {
	bot, repo := newTestBot(t)
	all, _ := repo.GetAll(context.Background())

	response, ok := bot.reply(context.Background(), "!tarot please", "Viewer")
	if !ok {
		t.Fatal("expected a reply")
	}
	if !strings.HasPrefix(response, "@Viewer ") {
		t.Errorf("expected reply to mention the viewer, got %q", response)
	}
	// With MockRandomizer returning 0, the draw is the first card, upright.
	for _, want := range []string{all[0].Name, "(Upright)", all[0].Meanings.Upright} {
		if !strings.Contains(response, want) {
			t.Errorf("expected %q in reply %q", want, response)
		}
	}
}

func TestReplySpread(t *testing.T) {
	bot, repo := newTestBot(t)
	all, _ := repo.GetAll(context.Background())

	response, ok := bot.reply(context.Background(), "!spread", "Viewer")
	if !ok {
		t.Fatal("expected a reply")
	}
	if !strings.HasPrefix(response, "@Viewer ") {
		t.Errorf("expected reply to mention the viewer, got %q", response)
	}
	for i, position := range []string{"Past", "Present", "Future"} {
		want := position + ": " + all[i].Name
		if !strings.Contains(response, want) {
			t.Errorf("expected %q in reply %q", want, response)
		}
	}
}

func TestReplyIgnoresOtherMessages(t *testing.T) {
	bot, _ := newTestBot(t)

	for _, message := range []string{"", "   ", "hello", "server", "!serverinfo", "!tarots", "join the !discord"} {
		if response, ok := bot.reply(context.Background(), message, "Viewer"); ok {
			t.Errorf("%q: expected no reply, got %q", message, response)
		}
	}
}
