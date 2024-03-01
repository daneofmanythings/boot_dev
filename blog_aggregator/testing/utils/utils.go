package utils

import (
	"github.com/daneofmanythings/blog_aggregator/testing/models"
)

func PostFeedsResponseToFeeds(pfr []models.PostFeedsResponse) []models.Feed {
	result := []models.Feed{}
	for _, p := range pfr {
		result = append(result, p.Feed)
	}
	return result
}
