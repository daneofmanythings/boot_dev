package models

import (
	"time"

	"github.com/google/uuid"
)

type Payload interface {
	pl()
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ApiKey    string
}

func (u User) pl() {}

type Feed struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	URL           string
	UserId        uuid.UUID
	LastUpdatedAt time.Time
}

func (u Feed) pl() {}

type FeedFollow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	FeedID    uuid.UUID
	UserID    uuid.UUID
}

func (u FeedFollow) pl() {}

type PostUserRequest struct {
	Name string
}

type PostFeedsRequest struct {
	Name string
	URL  string
}

type PostFeedsResponse struct {
	Feed       Feed
	FeedFollow FeedFollow `json:"feed_follow"`
}

type PostFeedFollowsRequest struct {
	FeedID uuid.UUID `json:"feed_id"`
}
