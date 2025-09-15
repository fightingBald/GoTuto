package main

import (
	"log"
	"net/http"

	appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/http"
	appsinmem "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inmem"
	app "github.com/fightingBald/GoTuto/apps/product-query-svc/app"
	httpadp "github.com/fightingBald/GoTuto/internal/adapters/http"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	repo := appsinmem.NewInMemRepo()
	svc := app.NewProductQueryService(repo)
	srv := appshttp.NewServer(svc)

	_ = httpadp.HandlerFromMux(srv, r)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
