package service

import (
	"context"
	"testing"
	"twitch-tarot-bot/internal/repository"
)

// MockRandomizer always returns 0 for deterministic testing
type MockRandomizer struct{}

func (m MockRandomizer) IntN(n int) int { return 0 }

func newTestService(t *testing.T) (*TarotService, *repository.CardRepository) {
	t.Helper()
	repo, err := repository.NewCardRepository()
	if err != nil {
		t.Fatalf("failed to load cards: %v", err)
	}
	return NewTarotService(repo, MockRandomizer{}), repo
}

func TestDraw(t *testing.T) {
	service, repo := newTestService(t)
	all, _ := repo.GetAll(context.Background())

	cards, err := service.Draw(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) != 3 {
		t.Fatalf("expected 3 cards, got %d", len(cards))
	}
	// With MockRandomizer returning 0, it should swap index 0 with 0, 1 with 1...
	// so the first 3 cards should exactly match the first 3 in the JSON.
	for i, card := range cards {
		if card.ID != all[i].ID {
			t.Errorf("card %d: expected %s, got %s", i, all[i].ID, card.ID)
		}
	}
}

func TestDrawRejectsInvalidCount(t *testing.T) {
	service, repo := newTestService(t)
	all, _ := repo.GetAll(context.Background())

	for _, n := range []int{-1, 0, len(all) + 1} {
		if _, err := service.Draw(context.Background(), n); err == nil {
			t.Errorf("expected error for n=%d", n)
		}
	}
}

func TestDrawDoesNotMutateDeck(t *testing.T) {
	repo, err := repository.NewCardRepository()
	if err != nil {
		t.Fatalf("failed to load cards: %v", err)
	}
	all, _ := repo.GetAll(context.Background())
	first := all[0].ID

	// lastRandomizer always picks the last remaining card, forcing a swap
	service := NewTarotService(repo, lastRandomizer{})
	if _, err := service.Draw(context.Background(), 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if all[0].ID != first {
		t.Fatalf("deck order changed: expected %s first, got %s", first, all[0].ID)
	}
}

type lastRandomizer struct{}

func (l lastRandomizer) IntN(n int) int { return n - 1 }

func TestDrawSingle(t *testing.T) {
	service, _ := newTestService(t)

	draw, err := service.DrawSingle(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if draw.IsReversed {
		t.Error("expected upright card when randomizer returns 0")
	}
	if draw.Meaning != draw.Card.Meanings.Upright {
		t.Errorf("expected upright meaning, got %q", draw.Meaning)
	}
	if draw.Position != nil {
		t.Errorf("expected no position, got %q", *draw.Position)
	}
}

func TestDrawThreeCardSpread(t *testing.T) {
	service, _ := newTestService(t)

	spread, err := service.DrawThreeCardSpread(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spread) != 3 {
		t.Fatalf("expected 3 cards, got %d", len(spread))
	}

	for i, want := range []string{"Past", "Present", "Future"} {
		if spread[i].Position == nil || *spread[i].Position != want {
			t.Errorf("card %d: expected position %s", i, want)
		}
	}
}

func TestDrawReversed(t *testing.T) {
	repo, err := repository.NewCardRepository()
	if err != nil {
		t.Fatalf("failed to load cards: %v", err)
	}
	service := NewTarotService(repo, lastRandomizer{})

	draw, err := service.DrawSingle(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !draw.IsReversed {
		t.Error("expected reversed card when randomizer returns 1")
	}
	if draw.Meaning != draw.Card.Meanings.Reversed {
		t.Errorf("expected reversed meaning, got %q", draw.Meaning)
	}
}
