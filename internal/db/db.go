// Package db contains database related functionality.
package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/vendelin8/card-fun/pkg/deck"
	"github.com/vendelin8/card-fun/pkg/util"
)

const (
	// maxRetries is the number of retries for Redis transaction optimistic locking.
	maxRetries = 1000
	// detailsFmt is how deck details is stored where %s is the deck id.
	dataFmt = "data%s"
)

var (
	rdb *redis.Client

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

func init() {
	redisAddr := util.GetEnv("REDIS_URL", "localhost:6379")
	rdb = redis.NewClient(&redis.Options{Addr: redisAddr, Password: "", DB: 0})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("Error trying to ping redis: %v", err)
	}
}

// StoreDeck stores the given cards in a new deck.
func StoreDeck(ctx context.Context, d *deck.Deck) error {
	pipe := rdb.Pipeline()
	var deckID string
	res := int64(1)
	for res > 0 {
		u, err := uuid.NewRandom()
		if err = handleError("create random uuid", err); err != nil {
			return err
		}
		deckID = u.String()
		res, err = pipe.Exists(ctx, deckID).Result()
		if err = handleError("check deck existance", err); err != nil {
			return err
		}
	}
	ifCodes := make([]interface{}, len(d.Codes))
	for i, c := range d.Codes {
		ifCodes[i] = c
	}
	pipe.RPush(ctx, deckID, ifCodes...)
	b, err := json.Marshal(d.Data)
	if err = handleError("JSON encode data", err); err != nil {
		return err
	}
	pipe.Set(ctx, fmt.Sprintf(dataFmt, deckID), string(b), 0)
	_, err = pipe.Exec(ctx)
	if err = handleError("create deck", err); err != nil {
		return err
	}
	d.ID = deckID
	return nil
}

// All returns a deck of cards for the given deck.
func All(ctx context.Context, d *deck.Deck) error {
	codes, err := rdb.LRange(ctx, d.ID, 0, -1).Result()
	if err = handleError("read deck", err); err != nil {
		return err
	}
	if len(codes) == 0 {
		return ErrMissingDeck
	}
	dataS, err := rdb.Get(ctx, fmt.Sprintf(dataFmt, d.ID)).Result()
	if err = handleError("read deck data", err); err != nil {
		return err
	}
	var data deck.Data
	if err = handleError("JSON decode deck data", json.Unmarshal([]byte(dataS), &data)); err != nil {
		return err
	}
	d.Data = data
	d.Codes = codes
	d.Remaining = len(codes)
	return nil
}

// Draw draws the top n cards of a deck for the given deck.
func Draw(ctx context.Context, d *deck.Deck, n int64) error {
	var lRangeRes *redis.StringSliceCmd
	txf := func(tx *redis.Tx) error {
		l, err := tx.LLen(ctx, d.ID).Result()
		if err = handleError("get deck length", err); err != nil {
			return err
		}
		if l == 0 {
			return ErrMissingDeck
		}
		if l < n {
			return ErrDeckLen{Remaining: int(l), Requested: int(n)}
		}

		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			lRangeRes = pipe.LRange(ctx, d.ID, 0, n-1)
			pipe.LTrim(ctx, d.ID, n, -1)
			if l == n {
				pipe.Del(ctx, fmt.Sprintf(dataFmt, d.ID))
			}
			return nil
		})
		return handleError("draw multiple cards", err)
	}

	// Retry if the key has been changed.
	for i := 0; i < maxRetries; i++ {
		err := rdb.Watch(ctx, txf, d.ID)
		if err == nil {
			codes, err := lRangeRes.Result()
			if err == nil {
				d.Codes = codes
				return nil
			}
		}
		if errors.Is(err, redis.TxFailedErr) {
			continue // optimistic lock lost, retry
		}
		return err
	}

	return ErrMaxRetries
}

func handleError(format string, err error, v ...any) error {
	if err == nil {
		return nil
	}
	format = fmt.Sprintf("%%v: failed to %s", format)
	if len(v) == 0 {
		log.Print(format)
		return util.ErrUnknown
	}
	v2 := make([]any, len(v)+1)
	v2[0] = err
	copy(v2[1:], v)
	log.Printf(format, v2...)
	return util.ErrUnknown
}
