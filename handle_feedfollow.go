package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vosram/aggregator/internal/database"
)

func handleFeedFollow(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: follow <url>")
	}
	url := cmd.Args[0]
	now := time.Now().UTC()
	currentUser, err := s.db.GetUser(context.Background(), s.conf.CurrentUser)
	if err != nil {
		return fmt.Errorf("couldn't fetch current user: %w", err)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	feedAdded, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create new feed follow: %w", err)
	}

	fmt.Println("Added Follow:")
	fmt.Println("User: ", feedAdded.UserName)
	fmt.Println("Feed: ", feedAdded.FeedName)
	return nil
}
