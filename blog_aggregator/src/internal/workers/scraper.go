package workers

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/daneofmanythings/blog_aggregator/internal/config"
	"github.com/daneofmanythings/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func fetchFeedData(url string) RSS {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Could not fetchFeedData: %s", err.Error())
	}
	rssFeed := RSS{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rssFeed)
	if err != nil {
		log.Fatalf("Could not Unmarshal xml: %s", err.Error())
	}
	return rssFeed
}

func processRSSFeed(app *config.Config, rss RSS, feedID uuid.UUID) {
	for _, post := range rss.Channel.Item {
		_, err := app.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       post.Title,
			Url:         post.Link,
			Description: validateNullString(post.Description),
			PublishedAt: validateNullTime(post.PubDate),
			FeedID:      feedID,
		})
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				continue
			}
			log.Fatalf("Unable to CreatePost: %s", err.Error())
		}
	}
}

func validateNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func validateNullTime(d string) sql.NullTime {
	date, err := detectFormatAndReturnDate(d)
	if err != nil {
		log.Println(err)
	}

	return sql.NullTime{
		Time:  date,
		Valid: err == nil,
	}
}

func detectFormatAndReturnDate(date string) (time.Time, error) {
	dateFormats := []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
	}
	for _, dateFormat := range dateFormats {
		date, err := time.Parse(dateFormat, date)
		if err != nil {
			continue
		}
		return date, nil
	}
	return time.Unix(0, 0), errors.New(fmt.Sprintf("Unable to parse date string: %s", date))
}

type RSSInfo struct {
	Url    string
	FeedID uuid.UUID
}

func FeedFetchingWorker(app *config.Config) {
	ticker := time.NewTicker(app.ScraperInterval)
	for ; ; <-ticker.C {
		feeds, err := app.DB.GetNextFeedsToFetch(context.Background(), app.NumFeedsToScrape)
		if err != nil {
			log.Fatalf("Could not get feeds: %s ( FeedFetchingWorker)", err.Error())
		}

		// WARN: this is kinda sus. maybe refactor and just pass the whole feed to the processRSSFeed func
		urls := []RSSInfo{}
		for _, feed := range feeds {
			urls = append(urls, RSSInfo{
				Url:    feed.Url,
				FeedID: feed.ID,
			})
			// log.Printf("Marking feed as fetched: %s", feed.ID)
			err = app.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
				ID: feed.ID,
				LastFetchedAt: sql.NullTime{
					Time:  time.Now().UTC(),
					Valid: true,
				},
			})
			if err != nil {
				log.Fatalf("Could not MarkFeedFetched: %s", err.Error())
			}
		}

		var wg sync.WaitGroup
		for _, url := range urls {
			wg.Add(1)
			go func(rssInfo RSSInfo) {
				defer wg.Done()
				rss := fetchFeedData(rssInfo.Url)
				processRSSFeed(app, rss, rssInfo.FeedID)
			}(url)
		}
		wg.Wait()
	}
}
