// Package db contains database related functionality.
package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/vendelin8/card-fun/pkg/deck"
)

var (
	ErrMissingDeck = errors.New("This deck cannot be found")
	ErrMaxRetries  = errors.New("Reached maximum number of retries")
)

type ErrDeckLen struct {
	Remaining int
	Requested int
}

func (e ErrDeckLen) Error() string {
	return fmt.Sprintf("Deck has less (%d) cards then requested (%d)", e.Remaining, e.Requested)
}

// StoreDeck stores the given cards in a new deck.
func StoreDeck(ctx context.Context, d *deck.Deck) error {
	fmt.Println("db.StoreDeck") // TODO: implement
	return nil
}

// All returns a deck of cards for the given deck.
func All(ctx context.Context, d *deck.Deck) error {
	fmt.Println("db.All") // TODO: implement
	return nil
}

// Draw draws the top n cards of a deck for the given deck.
func Draw(ctx context.Context, d *deck.Deck, n int64) error {
	fmt.Println("db.Draw") // TODO: implement
	return nil
}
