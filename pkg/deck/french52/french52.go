// Package french52 contains functionality for 52 card French deck.
package french52

import (
	"github.com/vendelin8/card-fun/pkg/deck"
)

type French52 deck.Deck

// New creates a new deck from given codes, or full deck when it's missing.
// It shuffles the deck if shuffled is true.
func (f *French52) New() error {
	return nil
}

// Resolve sets Cards property based on Codes property. No error checking,
// all codes must be valid.
func (f *French52) Resolve() {
}
