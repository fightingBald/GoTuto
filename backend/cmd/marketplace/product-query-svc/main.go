package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	appsinmem "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/inmem"
	appspg "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/postgres"
	productapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/product"
	userapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/user"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dsnFlag := flag.String("db-dsn", "", "Postgres DSN (if empty, use in-memory repo)")
	flag.Parse()

	// 支持 env 回退
	dsn := *dsnFlag
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	// address env fallback
	if *addr == ":8080" { // only override default
		if envAddr := os.Getenv("HTTP_ADDRESS"); envAddr != "" {
			*addr = envAddr
		}
	}

	log.Println("starting product-query-svc")

	var (
		repo     outbound.ProductRepository
		userRepo outbound.UserRepository
		pool     *pgxpool.Pool
	)

	// If DSN provided, use Postgres wiring
	if dsn != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		p, err := pgxpool.New(ctx, dsn)
		if err != nil {
			log.Fatalf("connect pg: %v", err)
			return
		}
		pool = p

		repo = appspg.NewProductRepository(pool)
		userRepo = appspg.NewUserRepository(pool)
	} else {
		store := appsinmem.NewInMemRepo()
		repo = store
		userRepo = store
	}

	// build service
	productSvc := productapp.NewService(repo)
	userSvc := userapp.NewService(userRepo)

	server := appshttp.NewServer(productSvc, userSvc)

	r := chi.NewRouter()
	// 注册健康检查
	r.HandleFunc("/healthz", server.Health)
	// 注册 OpenAPI 生成的 strict handler 到 chi Router
	strict := appshttp.NewStrictHTTPHandler(server, nil)
	h := appshttp.HandlerFromMux(strict, r)

	srv := &http.Server{
		Addr:    *addr,
		Handler: h,
	}

	// 启动服务器
	go func() {
		log.Printf("listening on %s", *addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctxShut, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShut); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	if pool != nil {
		pool.Close()
	}

	log.Println("server stopped")
}
