package routes

import (
	"net/http"

	"github.com/daneofmanythings/blog_aggregator/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func v1Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/readiness", handlers.Repo.GetReadinessHandler)
	r.Get("/err", handlers.Repo.GetErrHandler)

	r.Post("/users", handlers.Repo.V1PostUsersHandler)
	r.Get("/users", handlers.Repo.MiddlewareAuth(handlers.Repo.V1GetUsersHandler))

	r.Post("/feeds", handlers.Repo.MiddlewareAuth(handlers.Repo.V1PostFeedsHandler))
	r.Get("/feeds", handlers.Repo.V1GetFeedsHandler)

	r.Post("/feed_follows", handlers.Repo.MiddlewareAuth(handlers.Repo.V1PostFeedFollowsHandler))
	r.Delete("/feed_follows/{feedFollowID}", handlers.Repo.MiddlewareAuth(handlers.Repo.V1DeleteFeedFollowsHandler))
	r.Get("/feed_follows", handlers.Repo.MiddlewareAuth(handlers.Repo.V1GetFeedFollowsHandler))

	r.Get("/posts", handlers.Repo.MiddlewareAuth(handlers.Repo.V1GetPostsByUserHandler))

	return r
}
