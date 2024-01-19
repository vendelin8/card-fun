// Package testutil contains utility functions for testing.
package testutil

import (
	"testing"

	"github.com/vendelin8/card-fun/pkg/util"
)

const wantStatus = 400

// CompareErr checks if the given errors match.
func CompareErr(t *testing.T, wantResp string, resp *util.Err) {
	t.Helper()
	if wantResp != resp.Msg {
		t.Fatalf("response should be '%s', but it's: '%s'", wantResp, resp.Msg)
	}
	if wantStatus != resp.Code {
		t.Fatalf("status should be %d, but it's: %d", wantStatus, resp.Code)
	}
}

// NilErr checks if the given value is an error.
func NilErr(t *testing.T, err error, descr string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s error should be nil, but it's: %v", descr, err)
	}
}
