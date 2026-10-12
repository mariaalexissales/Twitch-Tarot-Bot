package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"twitch-tarot-bot/internal/domain"
	"twitch-tarot-bot/internal/testutil"
)

// fakeDeck is an in-memory CardSource
type fakeDeck struct {
	cards []domain.Card
	err   error
}

func (f fakeDeck) GetAll(ctx context.Context) ([]domain.Card, error) {
	return f.cards, f.err
}

func newFakeDeck(n int) fakeDeck {
	cards := make([]domain.Card, n)
	for i := range cards {
		cards[i] = domain.Card{
			ID:   fmt.Sprintf("card_%d", i),
			Name: fmt.Sprintf("Card %d", i),
			Meanings: domain.Meanings{
				Upright:  fmt.Sprintf("upright %d", i),
				Reversed: fmt.Sprintf("reversed %d", i),
			},
		}
	}
	return fakeDeck{cards: cards}
}

func TestDraw(t *testing.T) {
	deck := newFakeDeck(5)
	service := NewTarotService(deck, testutil.FirstRandomizer{})

	cards, err := service.Draw(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) != 3 {
		t.Fatalf("expected 3 cards, got %d", len(cards))
	}
	// With FirstRandomizer, it should swap index 0 with 0, 1 with 1...
	// so the first 3 cards should exactly match the first 3 in the deck.
	for i, card := range cards {
		if card.ID != deck.cards[i].ID {
			t.Errorf("card %d: expected %s, got %s", i, deck.cards[i].ID, card.ID)
		}
	}
}

func TestDrawHasNoDuplicates(t *testing.T) {
	deck := newFakeDeck(5)
	service := NewTarotService(deck, testutil.LastRandomizer{})

	cards, err := service.Draw(context.Background(), len(deck.cards))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	seen := make(map[string]bool)
	for _, card := range cards {
		if seen[card.ID] {
			t.Errorf("card %s drawn twice", card.ID)
		}
		seen[card.ID] = true
	}
}

func TestDrawRejectsInvalidCount(t *testing.T) {
	deck := newFakeDeck(5)
	service := NewTarotService(deck, testutil.FirstRandomizer{})

	for _, n := range []int{-1, 0, len(deck.cards) + 1} {
		if _, err := service.Draw(context.Background(), n); err == nil {
			t.Errorf("expected error for n=%d", n)
		}
	}
}

func TestDrawReturnsSourceError(t *testing.T) {
	want := errors.New("deck unavailable")
	service := NewTarotService(fakeDeck{err: want}, testutil.FirstRandomizer{})

	if _, err := service.Draw(context.Background(), 1); !errors.Is(err, want) {
		t.Errorf("expected %v, got %v", want, err)
	}
	if _, err := service.DrawSingle(context.Background()); !errors.Is(err, want) {
		t.Errorf("expected %v, got %v", want, err)
	}
	if _, err := service.DrawThreeCardSpread(context.Background()); !errors.Is(err, want) {
		t.Errorf("expected %v, got %v", want, err)
	}
}

func TestDrawDoesNotMutateDeck(t *testing.T) {
	deck := newFakeDeck(5)
	first := deck.cards[0].ID

	// LastRandomizer always picks the last remaining card, forcing a swap
	service := NewTarotService(deck, testutil.LastRandomizer{})
	if _, err := service.Draw(context.Background(), 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if deck.cards[0].ID != first {
		t.Fatalf("deck order changed: expected %s first, got %s", first, deck.cards[0].ID)
	}
}

func TestDrawSingle(t *testing.T) {
	service := NewTarotService(newFakeDeck(5), testutil.FirstRandomizer{})

	draw, err := service.DrawSingle(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if draw.IsReversed {
		t.Error("expected upright card when randomizer returns 0")
	}
	if draw.Orientation() != "Upright" {
		t.Errorf("expected Upright orientation, got %q", draw.Orientation())
	}
	if draw.Meaning != draw.Card.Meanings.Upright {
		t.Errorf("expected upright meaning, got %q", draw.Meaning)
	}
	if draw.Position != nil {
		t.Errorf("expected no position, got %q", *draw.Position)
	}
}

func TestDrawReversed(t *testing.T) {
	service := NewTarotService(newFakeDeck(5), testutil.LastRandomizer{})

	draw, err := service.DrawSingle(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !draw.IsReversed {
		t.Error("expected reversed card when randomizer returns 1")
	}
	if draw.Orientation() != "Reversed" {
		t.Errorf("expected Reversed orientation, got %q", draw.Orientation())
	}
	if draw.Meaning != draw.Card.Meanings.Reversed {
		t.Errorf("expected reversed meaning, got %q", draw.Meaning)
	}
}

func TestDrawThreeCardSpread(t *testing.T) {
	service := NewTarotService(newFakeDeck(5), testutil.FirstRandomizer{})

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

func TestDrawThreeCardSpreadNeedsThreeCards(t *testing.T) {
	service := NewTarotService(newFakeDeck(2), testutil.FirstRandomizer{})

	if _, err := service.DrawThreeCardSpread(context.Background()); err == nil {
		t.Error("expected error when the deck has fewer than 3 cards")
	}
}
