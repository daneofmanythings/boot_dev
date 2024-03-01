package runner

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/daneofmanythings/blog_aggregator/testing/models"
	"github.com/daneofmanythings/blog_aggregator/testing/utils"
)

func (r *Runner) RunPostUsersTests() {
	r.TestBegin("RunPostUsersTests")
	defer r.TestEnd()

	caller := "PostUsers"
	requestTemplate := RequestParameters{
		URL:           r.URL,
		Method:        http.MethodPost,
		Endpoint:      usersEndpoint,
		PathParameter: "",
		ContentType:   typeAppJSON,
		Headers:       map[string]string{},
		Body:          nil, // This will vary for the requests
	}

	names := []string{"Bob", "Tom", "Jan", "Ari"}
	for _, name := range names {
		requestTemplate.Body = models.PostUserRequest{
			Name: name,
		}
		payload, responseCode := SendRequest(
			requestTemplate,
			caller,
			r.Client)

		if responseCode != http.StatusCreated {
			utils.StatusCodeError(http.StatusCreated, responseCode, caller)
		}

		user := models.User{}
		err := json.Unmarshal(payload, &user)
		if err != nil {
			utils.JSONUnmarshalError(err, caller)
		}

		if user.Name != name {
			log.Fatalf("Recieved invalid payload data, Name: %s (%s)", user.Name, name)
		}
		r.Users = append(r.Users, user)
	}
}

func (r *Runner) RunGetUsersTests() {
	r.TestBegin("RunGetUsersTests")
	defer r.TestEnd()

	caller := "GetUsers"
	requestTemplate := RequestParameters{
		URL:           r.URL,
		Method:        http.MethodGet,
		Endpoint:      usersEndpoint,
		PathParameter: "",
		ContentType:   "",
		Headers: map[string]string{
			"Authorization": "", // This will vary for the requests
		},
		Body: nil,
	}
	for _, user := range r.Users {
		requestTemplate.Headers["Authorization"] = formApiAuthHeader(user.ApiKey)
		payload, responseCode := SendRequest(requestTemplate, caller, r.Client)

		if responseCode != http.StatusOK {
			utils.StatusCodeError(http.StatusOK, responseCode, caller)
		}
		userResponse := models.User{}
		err := json.Unmarshal(payload, &userResponse)
		if err != nil {
			utils.JSONUnmarshalError(err, caller)
		}

		if userResponse != user {
			log.Fatalf("Recieved invalid payload data, user: %s (%s)", userResponse, caller)
		}
	}
}
