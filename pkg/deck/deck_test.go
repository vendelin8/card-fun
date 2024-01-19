package deck_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/vendelin8/card-fun/pkg/deck"
	"github.com/vendelin8/card-fun/pkg/testutil"
)

func TestCheckDuplicatesOK(t *testing.T) {
	cases := []struct {
		value string
	}{
		{
			value: "a,b",
		},
		{
			value: "a,b,c,d",
		},
		{
			value: "a,b,c,d,1,2,3,4",
		},
	}

	for _, c := range cases {
		t.Run(c.value, func(t *testing.T) {
			err := deck.CheckDuplicates(strings.Split(c.value, ","))
			testutil.NilErr(t, err, "check duplicates")
		})
	}
}

func TestCheckDuplicatesErr(t *testing.T) {
	cases := []struct {
		value string
	}{
		{
			value: "a,b,a",
		},
		{
			value: "a,b,c, b",
		},
		{
			value: "a,b,c,d,1,2,3,4,1 ",
		},
	}

	want := deck.ErrDuplCard
	for _, c := range cases {
		t.Run(c.value, func(t *testing.T) {
			err := deck.CheckDuplicates(strings.Split(c.value, ","))
			if !errors.Is(err, want) {
				t.Fatalf("error should be %v, but it's: %v", want, err)
			}
		})
	}
}
