// Package deck contains base deck related structures and functions.
package deck

import (
	"errors"
	"fmt"
)

var ErrDuplCard = errors.New("Duplicate card found in input cards")

// ErrCardCodeLen is an error that shows a given card code with wrong len other than 2.
type ErrCardCodeLen struct {
	code    string
	wantLen int
}

func NewErrCardCodeLen(c string, e int) ErrCardCodeLen {
	return ErrCardCodeLen{code: c, wantLen: e}
}

func (e ErrCardCodeLen) Error() string {
	return fmt.Sprintf("Card code must be exactly %d caracters, but '%s' fails", e.wantLen, e.code)
}

// ErrCardCodeMiss is an error that shows a given card code is invalid because suit or
// value character is not matching the available cards of the deck.
type ErrCardCodeMiss struct {
	code string
}

func NewErrCardCodeMiss(c string) ErrCardCodeMiss {
	return ErrCardCodeMiss{code: c}
}

func (e ErrCardCodeMiss) Error() string {
	return fmt.Sprintf("Card code '%s' is NOT valid", e.code)
}

// DeckType is a blueprint for deck types and their possible actions.
type DeckType interface {
	New() error
	Resolve()
}

// Data contains everything that should be stored in db for deck id next to card codes.
type Data struct {
	Shuffled bool `json:"shuffled"`
}

// Details contains basic information for a deck.
type Details struct {
	Data
	ID        string `json:"deck_id"`
	Remaining int    `json:"remaining"`
}

// Deck contains all needed information of a deck.
type Deck struct {
	Details
	Codes []string            `json:"-"`
	Cards []map[string]string `json:"cards,omitempty"`
}

func (d *Deck) New() error {
	return nil
}

func (d *Deck) Resolve() {}
