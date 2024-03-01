package runner

import (
	"io"
	"log"
	"net/http"

	"github.com/daneofmanythings/blog_aggregator/testing/models"
	"github.com/google/uuid"
)

type Runner struct {
	// Manages all the data relating to tests
	Client *http.Client
	URL    string

	Users       []models.User
	Feeds       []models.Feed
	FeedFollows []models.FeedFollow
}

func (r *Runner) TestBegin(caller string) {
	log.Printf("=== %s ===\n", caller)
}

func (r *Runner) TestEnd() {
	log.Println("...ok")
}

func (r *Runner) GetUsersFeeds(user_id uuid.UUID) []models.Feed {
	feeds := []models.Feed{}
	for _, feed := range r.Feeds {
		if feed.UserId == user_id {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

func (r *Runner) GetUsersFeedFollows(user_id uuid.UUID) []models.FeedFollow {
	feedFollows := []models.FeedFollow{}
	for _, feedFollow := range r.FeedFollows {
		if feedFollow.UserID == user_id {
			feedFollows = append(feedFollows, feedFollow)
		}
	}
	return feedFollows
}

func (r *Runner) DeleteFeedFollow(ff_id uuid.UUID) {
	for i, feedFollow := range r.FeedFollows {
		if ff_id == feedFollow.ID {
			r.FeedFollows = append(r.FeedFollows[:i], r.FeedFollows[i+1:]...)
			return
		}
	}
}

func (r *Runner) AddFeedFollow(feedFollow models.FeedFollow) {
	r.FeedFollows = append(r.FeedFollows, feedFollow)
}

func (r *Runner) ResetDatabase(url string) {
	url = url + "/resetdatabase"
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %s", err.Error())
	}
	resp, err := r.Client.Do(req)
	if err != nil {
		log.Fatalf("Error resetting database: %s", err.Error())
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body for resetdatabase: %s", err.Error())
	}
	log.Println(string(body))
}
