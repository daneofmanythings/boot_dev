package runner

const (
	usersEndpoint       string = "/users"
	feedsEndpoint       string = "/feeds"
	feedFollowsEndpoint string = "/feed_follows"

	typeAppJSON string = "application/json"
)

func formApiAuthHeader(apikey string) string {
	return "ApiKey " + apikey
}
