package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	chi.RegisterMethod("LINK")
	chi.RegisterMethod("UNLINK")
	chi.RegisterMethod("WOOHOO")
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	http.ListenAndServe(":3333", r)
}
