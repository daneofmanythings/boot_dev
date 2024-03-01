package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/daneofmanythings/blog_aggregator/internal/database"
	"github.com/daneofmanythings/blog_aggregator/internal/handlers/utils"
	"github.com/google/uuid"
)

func (m *Repository) V1PostUsersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (decoding json payload)")
		return
	}
	user, err := m.App.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, user)
}

func (m *Repository) V1GetUsersHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	utils.RespondWithJSON(w, http.StatusOK, u)
}
