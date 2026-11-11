package mt5

import (
	"testing"

	"github.com/google/uuid"
)

func TestMT5Account_getHeaders(t *testing.T) {
	acc := &MT5Account{
		Id: uuid.Nil,
	}
	md := acc.getHeaders()
	if md == nil {
		t.Fatal("expected non-nil metadata")
	}
	val := md.Get("apikey")
	if len(val) == 0 || val[0] != "TRIAL" {
		t.Fatalf("expected apikey TRIAL, got %v", val)
	}

	acc.ApiKey = "custom-key"
	md = acc.getHeaders()
	val = md.Get("apikey")
	if len(val) == 0 || val[0] != "custom-key" {
		t.Fatalf("expected apikey custom-key, got %v", val)
	}
}
