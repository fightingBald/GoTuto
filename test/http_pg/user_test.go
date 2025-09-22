package http_pg_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	appspg "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/postgres"
	appsvc "github.com/fightingBald/GoTuto/apps/product-query-svc/app"
	"github.com/fightingBald/GoTuto/internal/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgconn"
)

func TestGetUserByID_Postgres(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := testutil.NewPool(ctx, t, pgDSN)
	defer pool.Close()
	if pgTemp {
		testutil.ApplyMigrations(ctx, t, pool)
	}

	insertStmt := `INSERT INTO users (name, email) VALUES ($1, $2)
		ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
		RETURNING id`

	var userID int64
	seedErr := pool.QueryRow(ctx, insertStmt, "Fixture User", "fixture@example.com").Scan(&userID)
	if seedErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(seedErr, &pgErr) && pgErr.Code == "42P01" {
			testutil.ApplyMigrations(ctx, t, pool)
			if err := pool.QueryRow(ctx, insertStmt, "Fixture User", "fixture@example.com").Scan(&userID); err != nil {
				t.Fatalf("seed user after migrations: %v", err)
			}
		} else {
			t.Fatalf("seed user: %v", seedErr)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
	})

	productRepo := appspg.NewProductRepository(pool)
	userRepo := appspg.NewUserRepository(pool)
	productSvc := appsvc.NewProductService(productRepo)
	userSvc := appsvc.NewUserService(userRepo)
	server := appshttp.NewServer(productSvc, userSvc)

	r := chi.NewRouter()
	h := appshttp.HandlerFromMux(server, r)

	ts := httptest.NewServer(h)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users/" + strconv.FormatInt(userID, 10))
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
	if user.Id == nil || *user.Id != userID {
		t.Fatalf("unexpected user id: %+v", user)
	}
	if user.Email == "" {
		t.Fatalf("expected email to be set: %+v", user)
	}
}
