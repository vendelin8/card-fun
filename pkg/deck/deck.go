// Package deck contains base deck related structures and functions.
package deck

import (
	"errors"
	"fmt"
	"strings"
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

// CheckDuplicates checks if a given list of card codes contain dulicates.
func CheckDuplicates(reqCodes []string) error {
	m := map[string]struct{}{}
	s := struct{}{}
	for i, c := range reqCodes {
		c = strings.TrimSpace(c)
		if _, ok := m[c]; ok {
			return ErrDuplCard
		}
		m[c] = s
		reqCodes[i] = c
	}
	return nil
}
