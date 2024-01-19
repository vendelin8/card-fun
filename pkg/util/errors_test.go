package util_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/vendelin8/card-fun/pkg/util"
	"github.com/vendelin8/card-fun/pkg/testutil"
)

func TestCheckErrNil(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	util.CheckErr(c, nil)
	b := c.Response().Body()
	if len(b) > 0 {
		t.Fatalf("ok should keep body empty, but it's: %s", b)
	}
}

func TestCheckErrErr(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	util.CheckErr(c, errors.New("some error"))
	b := c.Response().Body()
	var resp util.Err
	testutil.NilErr(t, json.Unmarshal(b, &resp), "unmarshal")
	testutil.CompareErr(t, "some error", &resp)
}
