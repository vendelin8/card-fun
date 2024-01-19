package util_test

import (
	"os"
	"testing"

	"github.com/vendelin8/card-fun/pkg/util"
	"github.com/vendelin8/card-fun/pkg/testutil"
)

func TestEnv(t *testing.T) {
	wantDef := "some default"
	key := "NOTEXISTS"
	res := util.GetEnv(key, wantDef)
	if res != wantDef {
		t.Fatalf("value should be %s, but it's: %s", wantDef, res)
	}

	value := "set value"
	testutil.NilErr(t, os.Setenv(key, value), "set env var")
	res = util.GetEnv(key, wantDef)
	if res != value {
		t.Fatalf("value should be %s, but it's: %s", value, res)
	}
}
