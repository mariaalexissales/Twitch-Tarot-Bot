package domain

import (
	"math/rand/v2"
)

// Randomizer allows us to inject deterministic randomness in tests.
type Randomizer interface {
	IntN(n int) int
}

// CryptoRandomizer uses math/rand/v2 which is safe for concurrent use.
type CryptoRandomizer struct{}

func (c CryptoRandomizer) IntN(n int) int {
	return rand.IntN(n)
}
