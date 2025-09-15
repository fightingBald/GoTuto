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

	appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/http"
	appsinmem "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inmem"
	appspg "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/postgres"
	"github.com/fightingBald/GoTuto/internal"
	httpadp "github.com/fightingBald/GoTuto/internal/adapters/http"
	"github.com/fightingBald/GoTuto/internal/ports"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dsnFlag := flag.String("db-dsn", "", "Postgres DSN (if empty, use in-memory repo)")
	migrate := flag.Bool("migrate", false, "run embedded SQL migrations when using Postgres")
	flag.Parse()

	// 支持 env 回退
	dsn := *dsnFlag
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}

	log.Println("starting product-query-svc")

	var (
		repo ports.ProductRepo
		pool *pgxpool.Pool
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

		if *migrate {
			if err := appspg.RunMigrations(context.Background(), pool); err != nil {
				log.Fatalf("run migrations: %v", err)
			}
		}

		repo = appspg.NewProductRepository(pool)
	} else {
		repo = appsinmem.NewInMemRepo()
	}

	// build service
	svc := internal.NewProductService(repo)

	server := appshttp.NewServer(svc)

	r := chi.NewRouter()
	// 注册健康检查
	r.HandleFunc("/healthz", server.Health)
	// 注册 OpenAPI 生成的 handler 到 chi Router
	h := httpadp.HandlerFromMux(server, r)

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
