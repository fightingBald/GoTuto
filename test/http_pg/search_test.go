package http_pg_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	appspg "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/postgres"
	commentapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/comment"
	productapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/product"
	userapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/user"
	"github.com/fightingBald/GoTuto/internal/testutil"
)

// TestSearchProducts_Postgres seeds are applied via migrations in dev/CI.
// This test requires DATABASE_URL to be set; otherwise it is skipped.
func TestSearchProducts_Postgres(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := testutil.NewPool(ctx, t, pgDSN)
	defer pool.Close()
	if pgTemp {
		testutil.ApplyMigrations(ctx, t, pool)
	}

	productRepo := appspg.NewProductRepository(pool)
	userRepo := appspg.NewUserRepository(pool)
	productSvc := productapp.NewService(productRepo)
	userSvc := userapp.NewService(userRepo)
	commentRepo := appspg.NewCommentRepository(pool)
	commentSvc := commentapp.NewService(commentRepo, productRepo, userRepo)
	server := appshttp.NewServer(productSvc, userSvc, commentSvc)

	h, err := appshttp.NewAPIHandler(server, nil)
	if err != nil {
		t.Fatalf("new api handler: %v", err)
	}

	ts := httptest.NewServer(h)
	defer ts.Close()

	t.Run("search wid returns 200", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/products/search?q=wid&page=1&pageSize=10", nil)
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("http do: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("unexpected status: %d", resp.StatusCode)
		}
		var pl appshttp.ProductList
		if err := json.NewDecoder(resp.Body).Decode(&pl); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if pl.Total < len(pl.Items) {
			t.Fatalf("expected total >= items length; got total=%d items=%d", pl.Total, len(pl.Items))
		}
	})
}
