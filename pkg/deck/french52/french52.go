// Package french52 contains functionality for 52 card French deck.
package french52

import (
	"math/rand"

	"github.com/vendelin8/card-fun/pkg/deck"
)

type French52 deck.Deck

var (
	valueMap = map[byte]string{
		'A': "ACE",
		'K': "KING",
		'Q': "QUEEN",
		'J': "JACK",
		'T': "10",
	}

	values = []byte("A23456789TJQK")

	suitMap = map[byte]string{
		'C': "CLUBS",
		'D': "DIAMONDS",
		'H': "HEARTS",
		'S': "SPADES",
	}
	suits = []byte("SDCH")
)

// New creates a new deck from given codes, or full deck when it's missing.
// It shuffles the deck if shuffled is true.
func (f *French52) New() error {
	if len(f.Codes) > 0 {
		if err := f.newFromCodes(); err != nil {
			return err
		}
	} else {
		f.newFull()
	}
	if f.Shuffled {
		codes := f.Codes
		rand.Shuffle(len(codes), func(i, j int) {
			codes[i], codes[j] = codes[j], codes[i]
		})
	}
	f.Remaining = len(f.Codes)
	return nil
}

// newFromCodes creates a deck from the given codes, may return a wrong input error.
func (f *French52) newFromCodes() (err error) {
	if err = deck.CheckDuplicates(f.Codes); err != nil {
		return
	}
	for _, c := range f.Codes {
		if len(c) != 2 {
			return deck.NewErrCardCodeLen(c, 2)
		}
		if len(resolveValue(c[0])) == 0 || len(resolveSuit(c[1])) == 0 {
			return deck.NewErrCardCodeMiss(c)
		}
	}
	return nil
}

// newFull creates a new full deck.
func (f *French52) newFull() {
	codes := make([]string, len(values)*len(suits))
	idx := 0
	code := make([]byte, 2)
	for _, s := range suits {
		code[1] = s
		for _, v := range values {
			code[0] = v
			codes[idx] = string(code)
			idx++
		}
	}
	f.Codes = codes
}

// Resolve sets Cards property based on Codes property. No error checking,
// all codes must be valid.
func (f *French52) Resolve() {
	l := len(f.Codes)
	cards := make([]map[string]string, l)
	for i, c := range f.Codes {
		cards[i] = map[string]string{
			"value": resolveValue(c[0]),
			"suit":  resolveSuit(c[1]),
			"code":  c,
		}
	}
	f.Cards = cards
}

// resolveValue resolves a card value (A, 2, 3, ..., T, J, Q, K) from a byte.
// returns empty string if invalid.
func resolveValue(c byte) string {
	if c >= '2' && c <= '9' {
		return string([]byte{c})
	}
	if v, ok := valueMap[c]; ok {
		return v
	}
	return ""
}

// resolveValue resolves a card suit (S, H, D, C) from a byte.
// returns empty string if invalid.
func resolveSuit(c byte) string {
	if v, ok := suitMap[c]; ok {
		return v
	}
	return ""
}
