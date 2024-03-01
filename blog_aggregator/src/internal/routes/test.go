package routes

import (
	"net/http"

	"github.com/daneofmanythings/blog_aggregator/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func testRouter() http.Handler {
	r := chi.NewRouter()

	// r.Get("/lastFetchedAt", handlers.Repo.TestGetNextFeedsToFetch)
	// r.Put("/markFeedFetched", handlers.Repo.TestMarkFeedFetched)
	r.Delete("/resetdatabase", handlers.Repo.ResetDatabase)

	return r
}
