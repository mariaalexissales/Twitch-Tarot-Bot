package repository

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"twitch-tarot-bot/internal/domain"
)

type CardRepository struct {
	cards []domain.Card
	rng   *rand.Rand
}

func NewCardRepository(jsonFilePath string) (*CardRepository, error) {
	fileBytes, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cards file: %w", err)
	}

	var cards []domain.Card
	if err := json.Unmarshal(fileBytes, &cards); err != nil {
		return nil, fmt.Errorf("failed to parse cards json: %w", err)
	}

	// Initialize thread-safe local randomizer
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	return &CardRepository{
		cards: cards,
		rng:   rng,
	}, nil
}

func (r *CardRepository) GetRandomCards(n int) []domain.Card {
	if n > len(r.cards) {
		n = len(r.cards)
	}

	perm := r.rng.Perm(len(r.cards))
	drawn := make([]domain.Card, n)

	for i := 0; i < n; i++ {
		drawn[i] = r.cards[perm[i]]
	}

	return drawn
}
