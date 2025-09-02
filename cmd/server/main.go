package main

import (
	mapi "github.com/fightingBald/GoTuto/ports/marketplaceapi" // 就是 module名 +/ports/marketplaceapi
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	s := mapi.NewServer()

	// 根据生成物注册路由
	_ = mapi.HandlerFromMux(s, r)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
