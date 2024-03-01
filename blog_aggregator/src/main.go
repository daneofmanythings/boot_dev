package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/daneofmanythings/blog_aggregator/internal/config"
	"github.com/daneofmanythings/blog_aggregator/internal/database"
	"github.com/daneofmanythings/blog_aggregator/internal/handlers"
	"github.com/daneofmanythings/blog_aggregator/internal/routes"
	"github.com/daneofmanythings/blog_aggregator/internal/workers"
	"github.com/joho/godotenv"
)

var app config.Config

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env: %s", err.Error())
	}

	portNumber := os.Getenv("PORT")
	dbURL := os.Getenv("DBURLMAIN")

	routers := []routes.RouterPath{
		routes.V1ROUTER,
	}

	debugMode := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debugMode {
		log.Println("Starting in DEBUG mode")
		routers = append(routers, routes.TESTROUTER)
		dbURL = os.Getenv("DBURLTEST")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error creating database: %s", err)
	}

	app.DB = database.New(db)
	app.ScraperInterval = time.Second * 60
	app.NumFeedsToScrape = 3

	repo := handlers.NewRepo(&app)
	handlers.LinkRepository(repo)

	handler := routes.InitRouter(routers)

	server := http.Server{
		Addr:    portNumber,
		Handler: handler,
	}

	go func() {
		workers.FeedFetchingWorker(repo.App)
	}()

	log.Printf("Starting server on port: %s", portNumber)
	log.Fatal(server.ListenAndServe())
}
