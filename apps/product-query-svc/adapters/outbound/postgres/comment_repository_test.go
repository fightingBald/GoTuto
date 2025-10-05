//go:build docker

package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupCommentRepo(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "app",
			"POSTGRES_PASSWORD": "app_password",
			"POSTGRES_DB":       "productdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	if err != nil {
		cancel()
		t.Fatalf("start container: %v", err)
	}

	cleanup := func() {
		cancel()
		_ = pgC.Terminate(context.Background())
	}

	host, err := pgC.Host(ctx)
	if err != nil {
		cleanup()
		t.Fatalf("host: %v", err)
	}
	port, err := pgC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		cleanup()
		t.Fatalf("mapped port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://app:app_password@%s:%s/productdb?sslmode=disable", host, port.Port())
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		cleanup()
		t.Fatalf("pgxpool.New: %v", err)
	}

	applyMigrations(ctx, pool, "migrations", t)

	return pool, func() {
		pool.Close()
		cleanup()
	}
}

func TestCommentRepository_WithDocker(t *testing.T) {
	if os.Getenv("SKIP_DOCKER_TESTS") == "1" {
		t.Skip("skipped by SKIP_DOCKER_TESTS=1")
	}

	pool, teardown := setupCommentRepo(t)
	defer teardown()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	productRepo := NewProductRepository(pool)
	commentRepo := NewCommentRepository(pool)

	product, err := domain.NewProduct("Fixture Gadget", 1299, nil)
	if err != nil {
		t.Fatalf("new product: %v", err)
	}
	productID, err := productRepo.Create(ctx, product)
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	var userID int64
	email := fmt.Sprintf("commenter-%d@example.com", time.Now().UnixNano())
	if err := pool.QueryRow(ctx, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", "Comment User", email).Scan(&userID); err != nil {
		t.Fatalf("insert user: %v", err)
	}

	first := &domain.Comment{ProductID: productID, UserID: userID, Content: "first"}
	firstID, err := commentRepo.CreateComment(ctx, first)
	if err != nil {
		t.Fatalf("create comment: %v", err)
	}
	if firstID == 0 {
		t.Fatalf("expected comment id to be set")
	}
	if first.CreatedAt.IsZero() || first.UpdatedAt.IsZero() {
		t.Fatalf("expected timestamps to be populated: %#v", first)
	}

	fetched, err := commentRepo.GetCommentByID(ctx, firstID)
	if err != nil {
		t.Fatalf("get comment: %v", err)
	}
	if fetched.Content != "first" || fetched.ProductID != productID || fetched.UserID != userID {
		t.Fatalf("unexpected fetched comment: %#v", fetched)
	}

	if err := first.UpdateContent("updated content"); err != nil {
		t.Fatalf("update content domain: %v", err)
	}
	if err := commentRepo.UpdateComment(ctx, first); err != nil {
		t.Fatalf("update comment: %v", err)
	}

	updated, err := commentRepo.GetCommentByID(ctx, firstID)
	if err != nil {
		t.Fatalf("get updated comment: %v", err)
	}
	if updated.Content != "updated content" {
		t.Fatalf("expected updated content, got %#v", updated)
	}

	second := &domain.Comment{ProductID: productID, UserID: userID, Content: "second"}
	secondID, err := commentRepo.CreateComment(ctx, second)
	if err != nil {
		t.Fatalf("create second comment: %v", err)
	}

	list, err := commentRepo.ListCommentsByProduct(ctx, productID)
	if err != nil {
		t.Fatalf("list comments: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(list))
	}
	if list[0].ID != secondID {
		t.Fatalf("expected newest comment first, got order %#v", list)
	}

	if err := commentRepo.DeleteComment(ctx, firstID); err != nil {
		t.Fatalf("delete comment: %v", err)
	}
	if _, err := commentRepo.GetCommentByID(ctx, firstID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}
