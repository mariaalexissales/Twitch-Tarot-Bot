package domain

import (
	"math/rand/v2"
)

// Randomizer allows us to inject deterministic randomness in tests.
type Randomizer interface {
	IntN(n int) int
}

// MathRandomizer draws from math/rand/v2's global source, which is randomly
// seeded and safe for concurrent use. It is not cryptographically secure.
type MathRandomizer struct{}

func (m MathRandomizer) IntN(n int) int {
	return rand.IntN(n)
}
