package handlers

import (
	"net/http"

	"github.com/daneofmanythings/blog_aggregator/internal/handlers/utils"
)

func (m *Repository) GetReadinessHandler(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	utils.RespondWithJSON(w, http.StatusOK, payload)
}

func (m *Repository) GetErrHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
