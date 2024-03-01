package handlers

import (
	"log"
	"net/http"

	"github.com/daneofmanythings/blog_aggregator/internal/handlers/utils"
)

func (m *Repository) ResetDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("ResetDatabase called...")
	err := m.App.DB.ResetUsers(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (reset users)")
	}
	err = m.App.DB.ResetFeeds(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (reset feeds)")
	}
	err = m.App.DB.ResetFeedFollows(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (reset feed follows)")
	}

	w.WriteHeader(http.StatusOK)
}
