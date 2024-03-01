package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/daneofmanythings/blog_aggregator/internal/database"
	"github.com/daneofmanythings/blog_aggregator/internal/handlers/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (m *Repository) V1PostFeedFollowsHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		FeedID string `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (decoding json payload)")
		return
	}

	feed_id, err := uuid.Parse(params.FeedID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error()+" (bad feed id)")
	}
	log.Println(feed_id)

	// WARN: repeat logic with feed creationg handler
	feed_follow, err := m.App.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed_id,
		UserID:    u.ID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (creating feed follow entry in database)")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, feed_follow)
}

func (m *Repository) V1DeleteFeedFollowsHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	feedFollowID := chi.URLParam(r, "feedFollowID")

	feedFollowUUID, err := uuid.Parse(feedFollowID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error()+" (inavalid feed follow uuid)")
		return
	}
	err = m.App.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowUUID,
		UserID: u.ID,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (deleting feed follow)")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Repository) V1GetFeedFollowsHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	feedFollows, err := m.App.DB.GetFeedFollows(r.Context(), u.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (getting feed follows for user)")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, feedFollows)
}
