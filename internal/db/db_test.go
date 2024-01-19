package db_test

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"testing"

	"github.com/vendelin8/card-fun/internal/db"
	"github.com/vendelin8/card-fun/pkg/deck"
	"github.com/vendelin8/card-fun/pkg/testutil"
)

var ctx = context.Background()

func TestDB(t *testing.T) {
	cases := []struct {
		codes    string
		shuffled bool
	}{
		{
			codes:    "AS",
			shuffled: true,
		},
		{
			codes:    "AS,2S,3S,4S,5S,6S,7S,8S,9S",
			shuffled: true,
		},
		{
			codes:    "AS,2H,3H,4H,5H,6H,7H,8H,9H",
			shuffled: true,
		},
		{
			codes: "AS,2H,3H,4H,5H,6H,7H,8H,9H",
		},
		{
			codes: "AS,2H,3H,4H,5H,6H,7H,8H,9H,AS,2C,3C,4C,5C,6C,7C,8C,9C",
		},
		{
			codes: "AS,2H,3H,4H,5H,6H,7H,8H,9H,TH,JH,QH,KH,AS,2C,3C,4C,5C,6C,7C,8C,9C,TC,JC,QC,KC",
		},
	}

	for _, c := range cases {
		ci := c
		t.Run(ci.codes, func(t *testing.T) {
			t.Parallel()
			codes := strings.Split(ci.codes, ",")
			d := &deck.Deck{Codes: codes}
			d.Shuffled = ci.shuffled
			testutil.NilErr(t, db.StoreDeck(ctx, d), "store deck")
			d.Codes = nil
			d.Shuffled = !d.Shuffled
			testutil.NilErr(t, db.All(ctx, d), "get all deck")
			allJoined := strings.Join(d.Codes, ",")
			if ci.codes != allJoined {
				t.Fatalf("original '%s' doesn't match ALL '%s'", ci.codes, allJoined)
			}
			if ci.shuffled != d.Shuffled {
				t.Fatalf("original shuffled %t doesn't match result %t", ci.shuffled, d.Shuffled)
			}
			rest := len(d.Codes)
			parts := []string{}
			for rest > 0 {
				n := rand.Intn(rest) + 1
				rest -= n
				testutil.NilErr(t, db.Draw(ctx, d, int64(n)), "draw from deck")
				parts = append(parts, strings.Join(d.Codes, ","))
				d.Codes = nil
			}
			partsJoined := strings.Join(parts, ",")
			if ci.codes != partsJoined {
				t.Fatalf("original codes '%s' doesn't match result '%s'", ci.codes, partsJoined)
			}
		})
	}
}

func TestMissingDeck(t *testing.T) {
	d := &deck.Deck{Details: deck.Details{ID: "someMissingDeck"}}
	err := db.All(ctx, d)
	if !errors.Is(err, db.ErrMissingDeck) {
		t.Fatalf("error should be %v, but it's: %v", db.ErrMissingDeck, err)
	}
}

func TestMissingSingleCard(t *testing.T) {
	d := &deck.Deck{Details: deck.Details{ID: "someMissingDeck"}}
	err := db.Draw(ctx, d, 1)
	if !errors.Is(err, db.ErrMissingDeck) {
		t.Fatalf("error should be %v, but it's: %v", db.ErrMissingDeck, err)
	}
}

func TestMissingMultiCard(t *testing.T) {
	d := &deck.Deck{Details: deck.Details{ID: "someMissingDeck"}}
	err := db.Draw(ctx, d, 2)
	if !errors.Is(err, db.ErrMissingDeck) {
		t.Fatalf("error should be %v, but it's: %v", db.ErrMissingDeck, err)
	}
}

func TestMissingMultiCard2(t *testing.T) {
	d := &deck.Deck{Codes: []string{"AS"}}
	testutil.NilErr(t, db.StoreDeck(ctx, d), "store last")

	err := db.Draw(ctx, d, 2)
	want := db.ErrDeckLen{1, 2}
	if !errors.Is(err, want) {
		t.Fatalf("error should be %v, but it's: %v", want, err)
	}
	testutil.NilErr(t, db.Draw(ctx, d, 1), "draw last")
	if len(d.Codes) != 1 || d.Codes[0] != "AS" {
		t.Fatalf("single draw should be 'AS' but it's '%s'", d.Codes)
	}

	err = db.Draw(ctx, d, 2)
	if !errors.Is(err, db.ErrMissingDeck) {
		t.Fatalf("error should be %v, but it's: %v", db.ErrMissingDeck, err)
	}
}
