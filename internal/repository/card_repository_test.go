package repository

import (
	"context"
	"testing"

	"twitch-tarot-bot/internal/domain"
)

func TestEmbeddedDeck(t *testing.T) {
	repo, err := NewCardRepository()
	if err != nil {
		t.Fatalf("failed to load cards: %v", err)
	}

	cards, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) != 78 {
		t.Fatalf("expected 78 cards, got %d", len(cards))
	}

	seen := make(map[string]bool)
	major := 0
	for _, card := range cards {
		if card.ID == "" || card.Name == "" {
			t.Errorf("card with empty id or name: %+v", card)
		}
		if seen[card.ID] {
			t.Errorf("duplicate card id %s", card.ID)
		}
		seen[card.ID] = true

		if card.Meanings.Upright == "" || card.Meanings.Reversed == "" {
			t.Errorf("%s: missing a meaning", card.ID)
		}

		switch card.ArcanaType {
		case domain.Major:
			major++
			if card.Suit != nil {
				t.Errorf("%s: major arcana should have no suit", card.ID)
			}
		case domain.Minor:
			if card.Suit == nil {
				t.Errorf("%s: minor arcana should have a suit", card.ID)
			}
		default:
			t.Errorf("%s: unknown arcana type %q", card.ID, card.ArcanaType)
		}
	}

	if major != 22 {
		t.Errorf("expected 22 major arcana, got %d", major)
	}
}
