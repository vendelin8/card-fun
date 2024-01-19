// Package db contains database related functionality.
package db

import (
	"context"
	"fmt"

	"github.com/vendelin8/card-fun/pkg/deck"
)

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
