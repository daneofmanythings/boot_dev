package config

import (
	"time"

	"github.com/daneofmanythings/blog_aggregator/internal/database"
)

type Config struct {
	DB               *database.Queries
	ScraperInterval  time.Duration
	NumFeedsToScrape int32
}
