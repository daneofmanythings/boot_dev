package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/daneofmanythings/chirpy/pkg/config"
	"gitlab.com/daneofmanythings/chirpy/pkg/handlers"
	"gitlab.com/daneofmanythings/chirpy/pkg/middleware"
)

const filepathRoot = "."

func Routes(app *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.MiddlewareCors)

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	fsHandler = middleware.MiddlewareMetricsInc(fsHandler, app)

	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	r_api := chi.NewRouter()
	r_api.Get("/healthz", handlers.Repo.ApiHealthHandler)
	r_api.Get("/reset", handlers.Repo.ApiResetMetrics)

	r_api.Post("/chirps", handlers.Repo.ApiPostChirpsHandler)
	r_api.Get("/chirps", handlers.Repo.ApiGetChirpsHandler)
	r_api.Get("/chirps/{chirpID}", handlers.Repo.ApiGetChirpsByChirpIDHandler)
	r_api.Delete("/chirps/{chirpID}", handlers.Repo.ApiDeleteChirpsByIDHandler)

	r_api.Post("/users", handlers.Repo.ApiPostUsersHandler)
	r_api.Put("/users", handlers.Repo.ApiPutUsersHandler)

	r_api.Post("/login", handlers.Repo.ApiPostLoginHandler)
	r_api.Post("/refresh", handlers.Repo.ApiRefreshTokenHandler)
	r_api.Post("/revoke", handlers.Repo.ApiRevokeTokenHandler)

	r_api.Post("/polka/webhooks", handlers.Repo.ApiPostPolkaWebhooksHandler)

	r.Mount("/api", r_api)

	r_admin := chi.NewRouter()
	r_admin.Get("/metrics", handlers.Repo.AdminHitsHandler)
	r.Mount("/admin", r_admin)

	return r
}
