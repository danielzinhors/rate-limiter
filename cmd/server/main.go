package main

import (
	"net/http"

	my_middleware "github.com/danielzinhors/rate-limiter/ratelimiter/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	rateLimiter := my_middleware.NewRateLimiter()

	r := chi.NewRouter()

	r.Use(rateLimiter)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
