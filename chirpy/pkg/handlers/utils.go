package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"gitlab.com/daneofmanythings/chirpy/pkg/models"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		log.Printf("Responding with 5xx error", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func censorChirp(body string) string {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	censorBar := "****"

	splitChirp := strings.Split(body, " ")
	for i, word := range splitChirp {
		if slices.Contains(badWords, strings.ToLower(word)) {
			splitChirp[i] = censorBar
		}
	}
	return strings.Join(splitChirp, " ")
}

func filterByAuthorID(chirps *[]models.Chirp, authorID int) []models.Chirp {
	filteredChirps := []models.Chirp{}
	for _, chirp := range *chirps {
		if chirp.AuthorID == authorID {
			filteredChirps = append(filteredChirps, chirp)
		}
	}
	return *&filteredChirps
}
