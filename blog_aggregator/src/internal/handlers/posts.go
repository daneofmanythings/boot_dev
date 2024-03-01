package handlers

import (
	"net/http"
	"strconv"

	"github.com/daneofmanythings/blog_aggregator/internal/database"
	"github.com/daneofmanythings/blog_aggregator/internal/handlers/utils"
)

func (m *Repository) V1GetPostsByUserHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	limStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limStr)
	if err != nil || limit < 1 {
		// utils.RespondWithError(w, http.StatusBadRequest, err.Error()+" ( Unable to convert limit to an integer)")
		// return
		limit = 5
	}

	posts, err := m.App.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: u.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" ( Unable to fetch feeds for user)")
		return
	}
	utils.RespondWithJSON(w, 200, posts)
}
