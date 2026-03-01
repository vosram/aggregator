package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vosram/aggregator/internal/database"
)

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: addfeed <feed name> <url>")
	}

	user, err := s.db.GetUser(context.Background(), s.conf.CurrentUser)
	if err != nil {
		return fmt.Errorf("couldn't get logged in user: %w", err)
	}

	now := time.Now().UTC()
	newFeed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating new feed: %w", err)
	}
	fmt.Println("New Feed Created:")
	fmt.Printf("- ID: %s\n", newFeed.ID)
	fmt.Printf("- CreatedAt: %v\n", newFeed.CreatedAt)
	fmt.Printf("- UpdatedAt: %v\n", newFeed.UpdatedAt)
	fmt.Printf("- Name: %s\n", newFeed.Name)
	fmt.Printf("- Url: %s\n", newFeed.Url)
	fmt.Printf("- UserID: %s\n", newFeed.UserID)
	return nil
}

func handleListFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get all feeds: %w", err)
	}
	fmt.Println("All Feeds")
	fmt.Println("=======")
	for i, feed := range feeds {
		fmt.Printf("Feed %d\n", i+1)
		fmt.Println("-------")
		fmt.Printf("- ID: %s\n", feed.ID)
		fmt.Printf("- CreatedAt: %s\n", feed.CreatedAt)
		fmt.Printf("- UpdatedAt: %v\n", feed.UpdatedAt)
		fmt.Printf("- Name: %s\n", feed.Name)
		fmt.Printf("- Url: %s\n", feed.Url)
		fmt.Printf("- UserID: %s\n", feed.UserID)
		fmt.Printf("- UserName: %s\n", feed.UserName)
		fmt.Println()
	}
	return nil
}
