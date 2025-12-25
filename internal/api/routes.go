package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	// Defining middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Defining Routes
	router.Get("/scrape", GetScrapedDataForProduct)

	return router
}
