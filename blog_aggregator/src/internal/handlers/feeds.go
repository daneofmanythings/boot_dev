package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/daneofmanythings/blog_aggregator/internal/database"
	"github.com/daneofmanythings/blog_aggregator/internal/handlers/utils"
	"github.com/google/uuid"
)

const duplicateURL string = "pq: duplicate key value violates unique constraint \"feeds_url_key\""

type parameters struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (m *Repository) V1PostFeedsHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (decoding json payload)")
		return
	}

	// getting the correct feed (it might already exist)
	feed := m.postFeedsHelper(w, r, &params, u)

	// WARN: repeat logic with feed follows handler
	feed_follow, err := m.App.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    u.ID,
	})
	if err != nil {
		if err.Error() == "pq: insert or update on table \"feed_follows\" violates foreign key constraint \"feed_follows_feed_id_fkey\"" {
			log.Println(feed.ID)
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (creating feed follow entry in database)")
		return
	}
	log.Println("feed follow created")
	log.Printf("feed_follow id: %s", feed_follow.ID)

	payload := struct {
		Feed       utils.Feed          `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{
		Feed:       utils.DatabaseFeedToFeed(&feed),
		FeedFollow: feed_follow,
	}

	utils.RespondWithJSON(w, http.StatusCreated, payload)
}

func (m *Repository) postFeedsHelper(w http.ResponseWriter, r *http.Request, params *parameters, u database.User) database.Feed {
	feed, err := m.App.DB.GetFeedFromURL(r.Context(), params.Url)
	if err != nil {
		feed, err := m.App.DB.CreateFeed(r.Context(), database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
			Url:       params.Url,
			UserID:    u.ID,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" (postFeedsHelper)")
		}
		log.Println("new feed created")
		return feed
	}
	log.Println("duplicate feed found")
	return feed
}

func (m *Repository) V1GetFeedsHandler(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := m.App.DB.GetAllFeeds(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error()+" ( retrieving feeds)")
	}
	feeds := []utils.Feed{}
	for _, feed := range dbFeeds {
		feeds = append(feeds, utils.DatabaseFeedToFeed(&feed))
	}
	utils.RespondWithJSON(w, http.StatusOK, feeds)
}
