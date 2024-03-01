package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/daneofmanythings/chirpy/internal/database"
	"gitlab.com/daneofmanythings/chirpy/pkg/config"
	"gitlab.com/daneofmanythings/chirpy/pkg/handlers"
	"gitlab.com/daneofmanythings/chirpy/routes"
)

const (
	portNumber   = ":8080"
	databasePath = "./database.json"
)

var app config.Config

func main() {
	debugMode := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debugMode {
		os.Remove(databasePath)
	}

	db, err := database.NewDB(databasePath)
	if err != nil {
		log.Fatalf("error creating database: %s", err)
	}

	app.ResetHits()
	app.DB = db

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env: %s", err.Error())
	}
	app.JWTSECRET = os.Getenv("JWT_SECRET")
	app.APIKEYPOLKA = os.Getenv("API_KEY_POLKA")

	repo := handlers.NewRepo(&app)
	handlers.LinkRepository(repo)

	server := http.Server{
		Addr:    portNumber,
		Handler: routes.Routes(&app),
	}

	log.Fatal(server.ListenAndServe())
}
