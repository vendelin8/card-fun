// Package deck contains base deck related structures and functions.
package deck

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
