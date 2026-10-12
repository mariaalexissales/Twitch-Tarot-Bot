// Package testutil holds helpers shared by tests.
package testutil

// FirstRandomizer always returns 0, so draws come off the top of the deck in
// order and every card is upright.
type FirstRandomizer struct{}

func (f FirstRandomizer) IntN(n int) int { return 0 }

// LastRandomizer always returns n-1, so draws come off the bottom of the deck
// and every card is reversed.
type LastRandomizer struct{}

func (l LastRandomizer) IntN(n int) int { return n - 1 }
