package service

import (
	"math/rand"
	"time"
	"twitch-tarot-bot/internal/domain"
	"twitch-tarot-bot/internal/repository"
)

type Tarot struct {
	repo *repository.CardRepository
	rng  *rand.Rand
}

func NewTarot(repo *repository.CardRepository) *Tarot {
	return &Tarot{
		repo: repo,
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func StringPtr(s string) *string {
	return &s
}

func (t *Tarot) DrawSingle() domain.DrawResult {
	cards := t.repo.GetRandomCards(1)
	return t.buildDrawResult(cards[0], nil)
}

func (t *Tarot) DrawThreeCardSpread() []domain.DrawResult {
	cards := t.repo.GetRandomCards(3)
	positions := []string{"Past", "Present", "Future"}

	results := make([]domain.DrawResult, 3)
	for i, card := range cards {
		results[i] = t.buildDrawResult(card, StringPtr(positions[i]))
	}
	return results
}

func (t *Tarot) buildDrawResult(card domain.Card, position *string) domain.DrawResult {
	meaning := card.Meanings.Upright
	isReversed := t.rng.Intn(2) == 1
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
