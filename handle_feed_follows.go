package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vosram/aggregator/internal/database"
)

func handleFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: follow <url>")
	}
	url := cmd.Args[0]
	now := time.Now().UTC()
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	feedAdded, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
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

func handleListFeedFollows(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("%s's feedFollows:\n", user.Name)
	fmt.Println("========")
	for _, feed := range feedFollows {
		fmt.Printf("* %s\n", feed.FeedName)
	}
	return nil
}
