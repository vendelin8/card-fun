package french52_test

import (
	"errors"
	"testing"

	"github.com/vendelin8/card-fun/pkg/deck"
	"github.com/vendelin8/card-fun/pkg/deck/french52"
	"github.com/vendelin8/card-fun/pkg/deck/french52/testhelp"
	decktestutil "github.com/vendelin8/card-fun/pkg/deck/testutil"
	"github.com/vendelin8/card-fun/pkg/testutil"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name      string
		reqCodes  []string
		wantCodes []string
		wantErr   error
	}{
		{
			name:     "from_req",
			reqCodes: testhelp.SomeCodes,
		},
		{
			name:     "from_req2",
			reqCodes: []string{"AS", "2S", "3S", "4C", "5H", "6D", "7D", "8S", "9H", "TD", "KC"},
		},
		{
			name:     "from_req_code_len_err",
			reqCodes: []string{"AS", "2S", "3S", "4C", "5H", "6D", "A"},
			wantErr:  deck.NewErrCardCodeLen("A", 2),
		},
		{
			name:     "from_req_duplicate",
			reqCodes: []string{"AS", "2S", "3S", "4C", "5H", "6D", "AS"},
			wantErr:  deck.ErrDuplCard,
		},
		{
			name:     "from_req_not_found_suit",
			reqCodes: []string{"AS", "2s", "3S", "4C", "5H", "6D"},
			wantErr:  deck.NewErrCardCodeMiss("2s"),
		},
		{
			name:     "from_req_not_found_value",
			reqCodes: []string{"AS", "1S", "3S", "4C", "5H", "6D"},
			wantErr:  deck.NewErrCardCodeMiss("1S"),
		},
		{
			name:      "full_deck",
			wantCodes: testhelp.AllCodes,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := french52.French52{Codes: c.reqCodes}
			err := d.New()
			if !errors.Is(err, c.wantErr) {
				t.Fatalf("new error should be %v, but it's: %v", c.wantErr, err)
			}
			if err != nil {
				return
			}
			if c.reqCodes == nil {
				c.reqCodes = c.wantCodes
			}
			if len(d.Codes) != len(c.reqCodes) {
				t.Fatalf("codes length should be %d, but it's: %d", len(c.reqCodes), len(d.Codes))
			}
			for i, dc := range d.Codes {
				if rc := c.reqCodes[i]; dc != rc {
					t.Fatalf("code at %d should be %s, but it's: %s", i, rc, dc)
				}
			}
			cCopy := make([]string, len(c.reqCodes)) // shuffling a copy
			copy(cCopy, c.reqCodes)
			e := french52.French52{Details: deck.Details{Data: deck.Data{Shuffled: true}}, Codes: cCopy}
			if err = e.New(); !errors.Is(err, c.wantErr) {
				t.Fatalf("new shuffled error should be %v, but it's: %v", c.wantErr, err)
			}
			if len(e.Codes) != len(d.Codes) {
				t.Fatalf("shuffled codes length should be %d, but it's: %d", len(e.Codes), len(d.Codes))
			}
			for i, ec := range e.Codes {
				if dc := d.Codes[i]; ec != dc {
					return
				}
			}
			t.Error("all shufled codes match")
		})
	}
}

func TestResolve(t *testing.T) {
	cases := []struct {
		name      string
		reqCodes  []string
		wantCards []map[string]string
	}{
		{
			name:      "from_req",
			reqCodes:  testhelp.SomeCodes,
			wantCards: testhelp.SomeCards,
		},
		{
			name:      "all",
			wantCards: testhelp.AllCards,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := french52.French52{Codes: c.reqCodes}
			testutil.NilErr(t, d.New(), "new french52 deck")
			d.Resolve()
			dd := (deck.Deck)(d)
			decktestutil.CompareDeck(t, len(c.wantCards), false, c.wantCards, &dd)
		})
	}
}
