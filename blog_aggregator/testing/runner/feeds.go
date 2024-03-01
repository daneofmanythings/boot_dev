package runner

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/daneofmanythings/blog_aggregator/testing/models"
	"github.com/daneofmanythings/blog_aggregator/testing/utils"
)

func (r *Runner) RunPostFeedsTest() {
	r.TestBegin("RunPostFeedsTests")
	defer r.TestEnd()

	caller := "PostFeeds"
	requestTemplate := RequestParameters{
		URL:           r.URL,
		Method:        http.MethodPost,
		Endpoint:      feedsEndpoint,
		PathParameter: "",
		ContentType:   typeAppJSON,
		Headers: map[string]string{ // api key required
			"Authorization": "",
		},
		Body: nil, // body required
	}
	for _, user := range r.Users {
		requestTemplate.Headers["Authorization"] = formApiAuthHeader(user.ApiKey)
		requestTemplate.Body = models.PostFeedsRequest{
			Name: "The Boot.dev Blog",
			URL:  "https://blog.boot.dev/index.xml",
		}

		payload, responseCode := SendRequest(
			requestTemplate,
			caller,
			r.Client)

		if responseCode != http.StatusCreated {
			utils.StatusCodeError(http.StatusCreated, responseCode, caller)
		}

		postFeedResp := models.PostFeedsResponse{}
		err := json.Unmarshal(payload, &postFeedResp)
		if err != nil {
			utils.JSONUnmarshalError(err, caller)
		}
		// if user.ID != postFeedResp.Feed.UserId || user.ID != postFeedResp.FeedFollow.UserID {
		// 	utils.WrongPayloadError(postFeedResp, caller)
		// }
		if !slices.Contains(r.Feeds, postFeedResp.Feed) {
			r.Feeds = append(r.Feeds, postFeedResp.Feed)
		}
		log.Println(postFeedResp.FeedFollow)
		r.FeedFollows = append(r.FeedFollows, postFeedResp.FeedFollow)
	}
}

func (r *Runner) RunGetFeedsTest() {
	r.TestBegin("RunGetFeedsTest")
	defer r.TestEnd()

	caller := "GetFeeds"
	requestTemplate := RequestParameters{
		URL:           r.URL,
		Method:        http.MethodGet,
		Endpoint:      feedsEndpoint,
		PathParameter: "",
		ContentType:   "",
		Headers:       map[string]string{},
		Body:          nil,
	}
	payload, responseCode := SendRequest(requestTemplate, caller, r.Client)

	if responseCode != http.StatusOK {
		utils.StatusCodeError(http.StatusOK, responseCode, caller)
	}

	feeds := []models.Feed{}
	err := json.Unmarshal(payload, &feeds)
	if err != nil {
		utils.JSONUnmarshalError(err, caller)
	}
	for _, feed := range feeds {
		if !slices.Contains(r.Feeds, feed) {
			utils.WrongPayloadError(feeds, caller)
		}
	}
}
