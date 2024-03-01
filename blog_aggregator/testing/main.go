package main

import (
	"log"
	"net/http"
	"os"

	"github.com/daneofmanythings/blog_aggregator/testing/runner"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env: %s", err.Error())
	}

	url := os.Getenv("URL")

	client := http.Client{}

	urlTest := url + "/test"
	urlV1 := url + "/v1"

	runner := runner.Runner{
		Client: &client,
		URL:    urlV1,
	}

	log.Printf("Testing on url: %s", urlV1)
	runner.ResetDatabase(urlTest)
	runner.RunPostUsersTests()
	runner.RunGetUsersTests()
	runner.RunPostFeedsTest()
	runner.RunGetFeedsTest()
	// runner.RunDeleteFeedFollowsTests()
	// runner.RunPostFeedFollowsTests()
	// runner.RunGetFeedFollowsTests()
}
