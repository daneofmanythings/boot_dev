package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

const filepathRoot = "."

type RouterPath string

const (
	V1ROUTER   RouterPath = "/v1"
	TESTROUTER RouterPath = "/test"
)

var routingPaths = map[RouterPath]routerFn{
	V1ROUTER:   v1Router,
	TESTROUTER: testRouter,
}

func routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	r.Handle("/app/*", fsHandler)
	r.Handle("/app", fsHandler)

	return r
}

type routerFn func() http.Handler

func InitRouter(routerPaths []RouterPath) http.Handler {
	mainRouter := routes()
	for _, routerPath := range routerPaths {
		mainRouter.Mount(string(routerPath), routingPaths[routerPath]())
	}
	return mainRouter
}
