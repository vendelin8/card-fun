package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/vendelin8/card-fun/internal/api"
	"github.com/vendelin8/card-fun/internal/db"
	"github.com/vendelin8/card-fun/pkg/util"
	"github.com/vendelin8/card-fun/pkg/deck"
	"github.com/vendelin8/card-fun/pkg/deck/french52/testhelp"
	decktestutil "github.com/vendelin8/card-fun/pkg/deck/testutil"
	"github.com/vendelin8/card-fun/pkg/testutil"
)

var (
	app    = fiber.New()
	reUUID = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
	dID    string
)

func init() {
	app.Put("/v1/create", api.CreateHandler)
	app.Get("/v1/open", api.OpenHandler)
	app.Post("/v1/draw", api.DrawHandler)
}

// call tests an API call to the given url with given method, want may be an error or
// a successful result to check.
func call(t *testing.T, url, method string, want interface{}) {
	t.Helper()
	req := httptest.NewRequest(method, url, nil)
	res, err := app.Test(req, -1)
	testutil.NilErr(t, err, "http call")
	defer res.Body.Close()

	wantErr, isErr := want.(error)
	var wantStatus int
	if isErr {
		wantStatus = 400
	} else {
		wantStatus = 200
	}

	if res.StatusCode != wantStatus {
		bs, _ := io.ReadAll(res.Body)
		t.Fatalf("http call status should be %d, but it's: %d with body %s",
			wantStatus, res.StatusCode, bs)
	}
	var buf bytes.Buffer
	tee := io.TeeReader(res.Body, &buf)
	dcd := json.NewDecoder(tee)

	if isErr {
		var e util.Err
		testutil.NilErr(t, dcd.Decode(&e), "JSON decode error")
		testutil.CompareErr(t, wantErr.Error(), &e)
		return
	}

	var result map[string]interface{}
	w := want.(map[string]interface{})
	testutil.NilErr(t, dcd.Decode(&result), "JSON decode")
	if len(result) != len(w) {
		t.Fatalf("result length should be %d, but it's: %#v", len(w), result)
	}

	var d deck.Deck
	dcd = json.NewDecoder(&buf)
	testutil.NilErr(t, dcd.Decode(&d), "JSON decode deck")
	if len(d.ID) == 0 {
		if _, ok := w["deck_id"]; ok {
			t.Fatal("deck id should be set, but it's NOT")
		}
	} else {
		wdID, ok := w["deck_id"]
		if !ok {
			t.Fatal("deck id should NOT be set, but it is")
		}
		deckID := wdID.(string)
		if len(deckID) == 0 {
			if !reUUID.MatchString(d.ID) {
				t.Fatalf("deck id should match uuid, but it's: %s", d.ID)
			}
		} else if deckID != d.ID {
			t.Fatalf("deck id should be %s, but it's: %s", deckID, d.ID)
		}
		dID = d.ID
	}
	remaining, _ := w["remaining"].(int)
	shuffled, _ := w["shuffled"].(bool)
	cards, _ := w["cards"].([]map[string]string)
	decktestutil.CompareDeck(t, remaining, shuffled, cards, &d)
}

// TestAPIFull tests a new full deck.
func TestAPIFull(t *testing.T) {
	call(t, "/v1/create", http.MethodPut, map[string]interface{}{
		"deck_id": "", "remaining": 52, "shuffled": false,
	})

	call(t, fmt.Sprintf("/v1/open?deck_id=%s", dID), http.MethodGet, map[string]interface{}{
		"deck_id": dID, "remaining": 52, "shuffled": false, "cards": testhelp.AllCards,
	})

	call(t, fmt.Sprintf("/v1/draw?deck_id=%s", dID), http.MethodPost,
		map[string]interface{}{"cards": testhelp.AllCards[:1]})

	call(t, fmt.Sprintf("/v1/draw?deck_id=%s&count=50", dID), http.MethodPost,
		map[string]interface{}{"cards": testhelp.AllCards[1:51]})

	call(t, fmt.Sprintf("/v1/draw?deck_id=%s&count=1", dID), http.MethodPost,
		map[string]interface{}{"cards": testhelp.AllCards[51:]})
}

// TestAPIFull tests a new partial deck from request.
func TestAPIReq(t *testing.T) {
	call(t, fmt.Sprintf("/v1/create?cards=%s", strings.Join(testhelp.SomeCodes, ",")),
		http.MethodPut, map[string]interface{}{
			"deck_id": "", "remaining": len(testhelp.SomeCodes), "shuffled": false,
		})

	call(t, fmt.Sprintf("/v1/open?deck_id=%s", dID), http.MethodGet, map[string]interface{}{
		"deck_id": "", "remaining": len(testhelp.SomeCodes), "shuffled": false,
		"cards": testhelp.SomeCards,
	})

	call(t, fmt.Sprintf("/v1/draw?deck_id=%s", dID), http.MethodPost,
		map[string]interface{}{"cards": testhelp.SomeCards[:1]})

	call(t, fmt.Sprintf("/v1/draw?deck_id=%s&count=10", dID), http.MethodPost,
		db.ErrDeckLen{Remaining: 9, Requested: 10})

	call(t, fmt.Sprintf("/v1/draw?deck_id=%s&count=9", dID), http.MethodPost,
		map[string]interface{}{"cards": testhelp.SomeCards[1:]})

	call(t, fmt.Sprintf("/v1/draw?deck_id=%s&count=1", dID), http.MethodPost,
		db.ErrMissingDeck)

	call(t, fmt.Sprintf("/v1/open?deck_id=%s", dID), http.MethodGet, db.ErrMissingDeck)
}
