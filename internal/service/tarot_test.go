package service

import (
	"context"
	"testing"
	"twitch-tarot-bot/internal/repository"
)

// MockRandomizer always returns 0 for deterministic testing
type MockRandomizer struct{}

func (m MockRandomizer) IntN(n int) int { return 0 }

func TestDrawThreeCardSpread(t *testing.T) {
	repo, _ := repository.NewCardRepository()
	service := NewTarotService(repo, MockRandomizer{})

	cards, err := service.Draw(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) != 3 {
		t.Fatalf("expected 3 cards, got %d", len(cards))
	}
	// With MockRandomizer returning 0, it should swap index 0 with 0, 1 with 1...
	// so the first 3 cards should exactly match the first 3 in the JSON.
}
