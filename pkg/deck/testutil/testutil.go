package testutil

import (
	"testing"

	"github.com/vendelin8/card-fun/pkg/deck"
)

func CompareDeck(t *testing.T, wantRem int, wantShuf bool, wantCards []map[string]string, res *deck.Deck) {
	t.Helper()
	if wantRem != res.Remaining {
		t.Fatalf("remaining should be %d, but it's: %d", wantRem, res.Remaining)
	}

	if wantShuf != res.Shuffled {
		t.Fatalf("shuffled should be %t, but it's: %t", wantShuf, res.Shuffled)
	}

	rCards := res.Cards
	if len(wantCards) != len(rCards) {
		t.Fatalf("result cards length should be %d, but it's: %d", len(wantCards), len(rCards))
	}
	for i, r := range rCards {
		w := wantCards[i]
		if len(r) != len(w) {
			t.Fatalf("result card at %d should key count mismatch, expected: %#v, but it's: %#v", i, w, r)
		}
		for j, ri := range r {
			wi := w[j]
			if ri != wi {
				t.Fatalf("result card key %s at %d should be %s, but it's: %s", j, i, wi, ri)
			}
		}
	}
}
