package service

import (
	"context"
	"fmt"

	"twitch-tarot-bot/internal/domain"
	"twitch-tarot-bot/internal/repository"
)

type TarotService struct {
	repo *repository.CardRepository
	rand domain.Randomizer
}

func NewTarotService(repo *repository.CardRepository, rand domain.Randomizer) *TarotService {
	return &TarotService{repo: repo, rand: rand}
}

func (s *TarotService) Draw(ctx context.Context, n int) ([]domain.Card, error) {
	cards, err := s.repo.GetAll(ctx)
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
