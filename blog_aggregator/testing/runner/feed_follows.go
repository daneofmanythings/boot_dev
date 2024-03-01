package runner

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/daneofmanythings/blog_aggregator/testing/models"
	"github.com/daneofmanythings/blog_aggregator/testing/utils"
)

func (r *Runner) RunDeleteFeedFollowsTests() {
	r.TestBegin("RunDeleteFeedFollowsTests")
	defer r.TestEnd()

	caller := "DeleteFeedFollows"
	requestTemplate := RequestParameters{
		URL:           r.URL,
		Method:        http.MethodDelete,
		Endpoint:      feedFollowsEndpoint,
		PathParameter: "", // This will vary. determines the follow to delete by id
		ContentType:   "",
		Headers: map[string]string{
			"Authorization": "", // requires auth
		},
		Body: nil,
	}
	for _, user := range r.Users {
		requestTemplate.Headers["Authorization"] = formApiAuthHeader(user.ApiKey)
		for _, feedFollow := range r.GetUsersFeedFollows(user.ID) {
			requestTemplate.PathParameter = feedFollow.ID.String()

			_, responseCode := SendRequest(requestTemplate, caller, r.Client)

			if responseCode != http.StatusOK {
				utils.StatusCodeError(http.StatusOK, responseCode, caller)
			}

			r.DeleteFeedFollow(feedFollow.ID)
		}
	}
}

func (r *Runner) RunPostFeedFollowsTests() {
	r.TestBegin("RunPostFeedFollowsTests")
	defer r.TestEnd()

	caller := "PostFeedFollows"
	requestTemplate := RequestParameters{
		URL:           r.URL,
		Method:        http.MethodPost,
		Endpoint:      feedFollowsEndpoint,
		PathParameter: "",
		ContentType:   typeAppJSON,
		Headers: map[string]string{
			"Authorization": "",
		},
		Body: nil,
	}
	for _, user := range r.Users {
		requestTemplate.Headers["Authorization"] = formApiAuthHeader(user.ApiKey)
		for _, feed := range r.GetUsersFeeds(user.ID) {
			requestTemplate.Body = models.PostFeedFollowsRequest{
				FeedID: feed.ID,
			}

			payload, responseCode := SendRequest(requestTemplate, caller, r.Client)

			if responseCode != http.StatusCreated {
				utils.StatusCodeError(http.StatusCreated, responseCode, caller)
			}

			feedFollow := models.FeedFollow{}
			err := json.Unmarshal(payload, &feedFollow)
			if err != nil {
				utils.JSONUnmarshalError(err, caller)
			}
			if feed.UserId != user.ID {
				utils.WrongPayloadError(feed, caller)
			}

			r.AddFeedFollow(feedFollow)
		}
	}
}

func (r *Runner) RunGetFeedFollowsTests() {
	r.TestBegin("RunGetFeedFollowsTests")
	defer r.TestEnd()

	caller := "GetFeedFollows"
	requestTemplate := RequestParameters{
		URL:           r.URL,
		Method:        http.MethodGet,
		Endpoint:      feedFollowsEndpoint,
		PathParameter: "",
		ContentType:   "",
		Headers: map[string]string{
			"Authorization": "", // requires authorization
		},
		Body: nil,
	}

	for _, user := range r.Users {
		requestTemplate.Headers["Authorization"] = formApiAuthHeader(user.ApiKey)
		payload, responseCode := SendRequest(requestTemplate, caller, r.Client)

		if responseCode != http.StatusOK {
			utils.StatusCodeError(http.StatusOK, responseCode, caller)
		}

		feedFollows := []models.FeedFollow{}
		err := json.Unmarshal(payload, &feedFollows)
		if err != nil {
			utils.JSONUnmarshalError(err, caller)
		}

		for _, feedFollow := range feedFollows {
			userFollows := r.GetUsersFeedFollows(user.ID)
			if !slices.Contains(userFollows, feedFollow) {
				utils.WrongPayloadError(feedFollows, caller)
			}
			if feedFollow.UserID != user.ID {
				log.Println("failed at user")
				utils.WrongPayloadError(feedFollows, caller)
			}
		}
	}
}
