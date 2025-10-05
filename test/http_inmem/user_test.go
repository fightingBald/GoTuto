package http_inmem_test

import (
	"encoding/json"
	"net/http"
	"testing"

	appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	appsinmem "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/inmem"
	"github.com/fightingBald/GoTuto/internal/testutil"
)

func TestGetUserByID_InMem(t *testing.T) {
	store := appsinmem.NewInMemRepo()
	ts := testutil.NewHTTPServer(store, store, store)
	t.Cleanup(ts.Close)

	resp, err := http.Get(ts.URL + "/users/1")
	if err != nil {
		t.Fatalf("http get user: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var user appshttp.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("decode user: %v", err)
	}
	if user.Id == nil || *user.Id != 1 {
		t.Fatalf("unexpected user id: %+v", user)
	}
	if user.Email != "alice@example.com" {
		t.Fatalf("unexpected user email: %+v", user)
	}

	resp404, err := http.Get(ts.URL + "/users/9999")
	if err != nil {
		t.Fatalf("http get user 404: %v", err)
	}
	defer resp404.Body.Close()
	if resp404.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp404.StatusCode)
	}
}
