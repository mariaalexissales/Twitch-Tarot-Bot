package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"twitch-tarot-bot/assets"
	"twitch-tarot-bot/internal/domain"
)

type CardRepository struct {
	cards []domain.Card
}

func NewCardRepository() (*CardRepository, error) {
	var cards []domain.Card
	if err := json.Unmarshal(assets.CardsJSON, &cards); err != nil {
		return nil, fmt.Errorf("failed to parse cards json: %w", err)
	}
	if len(cards) == 0 {
		return nil, errors.New("cards json contains no cards")
	}

	return &CardRepository{cards: cards}, nil
}

func (r *CardRepository) GetAll(ctx context.Context) ([]domain.Card, error) {
	return r.cards, nil
}
