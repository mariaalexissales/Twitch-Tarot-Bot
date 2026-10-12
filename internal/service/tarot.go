package service

import (
	"context"
	"fmt"

	"twitch-tarot-bot/internal/domain"
)

// CardSource supplies the deck to draw from.
type CardSource interface {
	GetAll(ctx context.Context) ([]domain.Card, error)
}

type TarotService struct {
	cards CardSource
	rand  domain.Randomizer
}

func NewTarotService(cards CardSource, rand domain.Randomizer) *TarotService {
	return &TarotService{cards: cards, rand: rand}
}

func (s *TarotService) Draw(ctx context.Context, n int) ([]domain.Card, error) {
	cards, err := s.cards.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if n <= 0 || n > len(cards) {
		return nil, fmt.Errorf("invalid number of cards requested: %d", n)
	}

	drawn := make([]domain.Card, len(cards))
	copy(drawn, cards)

	// Partial Fisher-Yates: shuffle only 'n' positions
	for i := 0; i < n; i++ {
		j := i + s.rand.IntN(len(drawn)-i)
		drawn[i], drawn[j] = drawn[j], drawn[i]
	}

	return drawn[:n], nil
}

func (s *TarotService) DrawSingle(ctx context.Context) (domain.DrawResult, error) {
	cards, err := s.Draw(ctx, 1)
	if err != nil {
		return domain.DrawResult{}, err
	}
	return s.buildDrawResult(cards[0], nil), nil
}

func (s *TarotService) DrawThreeCardSpread(ctx context.Context) ([]domain.DrawResult, error) {
	positions := []string{"Past", "Present", "Future"}

	cards, err := s.Draw(ctx, len(positions))
	if err != nil {
		return nil, err
	}

	results := make([]domain.DrawResult, len(cards))
	for i, card := range cards {
		results[i] = s.buildDrawResult(card, &positions[i])
	}
	return results, nil
}

func (s *TarotService) buildDrawResult(card domain.Card, position *string) domain.DrawResult {
	meaning := card.Meanings.Upright
	isReversed := s.rand.IntN(2) == 1
	if isReversed {
		meaning = card.Meanings.Reversed
	}

	return domain.DrawResult{
		Card:       card,
		IsReversed: isReversed,
		Meaning:    meaning,
		Position:   position,
	}
}
