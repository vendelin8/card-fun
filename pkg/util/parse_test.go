package util_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/vendelin8/card-fun/pkg/util"
	"github.com/vendelin8/card-fun/pkg/testutil"
)

type Response struct {
	Code int
	Msg  string
}

func TestParseBool(t *testing.T) {
	app := fiber.New()

	cases := []struct {
		value  string
		wantOK bool
		wantV  bool
	}{
		{
			value:  "wrong",
			wantOK: false,
			wantV:  false,
		},
		{
			value:  "not a bool false",
			wantOK: false,
			wantV:  false,
		},
		{
			value:  "true_neither_this_is_a_boolean",
			wantOK: false,
			wantV:  false,
		},
		{
			value:  "",
			wantOK: false,
			wantV:  false,
		},
		{
			value:  "0",
			wantOK: true,
			wantV:  false,
		},
		{
			value:  "1",
			wantOK: true,
			wantV:  true,
		},
		{
			value:  "2",
			wantOK: false,
			wantV:  false,
		},
		{
			value:  "true",
			wantOK: true,
			wantV:  true,
		},
		{
			value:  "True ",
			wantOK: true,
			wantV:  true,
		},
		{
			value:  "false",
			wantOK: true,
			wantV:  false,
		},
		{
			value:  " False ",
			wantOK: true,
			wantV:  false,
		},
	}

	var resp util.Err
	for _, c := range cases {
		t.Run(c.value, func(t *testing.T) {
			d := app.AcquireCtx(&fasthttp.RequestCtx{})
			v, ok := util.ParseBool(d, c.value, "awesome")
			if ok != c.wantOK {
				t.Fatalf("ok should be %t, but it's: %t", c.wantOK, ok)
			}
			if v != c.wantV {
				t.Fatalf("value should be %t, but it's: %t", c.wantV, v)
			}
			b := d.Response().Body()

			if !ok {
				testutil.NilErr(t, json.Unmarshal(b, &resp), "unmarshal")
				wantResp := fmt.Sprintf(
					"Parameter 'awesome' contains '%s' which is not a boolean value. Enter 'true' or 'false', or leave empty",
					c.value)
				testutil.CompareErr(t, wantResp, &resp)
			} else {
				if len(b) > 0 {
					t.Fatalf("ok should keep body empty, but it's: %s", b)
				}
			}
		})
	}
}

func TestParsePosInt(t *testing.T) {
	app := fiber.New()

	cases := []struct {
		value  string
		wantOK bool
		wantV  int64
	}{
		{
			value:  "wrong",
			wantOK: false,
			wantV:  0,
		},
		{
			value:  "not a bool 1",
			wantOK: false,
			wantV:  0,
		},
		{
			value:  "0_neither_this_is_an_integer",
			wantOK: false,
			wantV:  0,
		},
		{
			value:  "true",
			wantOK: false,
			wantV:  0,
		},
		{
			value:  "-1",
			wantOK: false,
			wantV:  0,
		},
		{
			value:  "0",
			wantOK: false,
			wantV:  0,
		},
		{
			value:  "1",
			wantOK: true,
			wantV:  1,
		},
		{
			value:  "2 ",
			wantOK: true,
			wantV:  2,
		},
		{
			value:  " 3",
			wantOK: true,
			wantV:  3,
		},
		{
			value:  " 4 ",
			wantOK: true,
			wantV:  4,
		},
	}

	var resp util.Err
	for _, c := range cases {
		t.Run(c.value, func(t *testing.T) {
			d := app.AcquireCtx(&fasthttp.RequestCtx{})
			v, ok := util.ParsePosInt(d, c.value, "awesome")
			if ok != c.wantOK {
				t.Fatalf("ok should be %t, but it's: %t", c.wantOK, ok)
			}
			if v != c.wantV {
				t.Fatalf("value should be %d, but it's: %d", c.wantV, v)
			}
			b := d.Response().Body()

			if !ok {
				testutil.NilErr(t, json.Unmarshal(b, &resp), "unmarshal")
				wantResp := fmt.Sprintf(
					"Parameter 'awesome' contains '%s' which is not a positive integer value. Enter eg. 2, or leave empty for 1",
					c.value)
				testutil.CompareErr(t, wantResp, &resp)
			} else {
				if len(b) > 0 {
					t.Fatalf("ok should keep body empty, but it's: %s", b)
				}
			}
		})
	}
}
