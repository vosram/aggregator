package main

import (
	"context"
	"fmt"
)

func handleFollowing(s *state, cmd command) error {

	currUser, err := s.db.GetUser(context.Background(), s.conf.CurrentUser)
	if err != nil {
		return fmt.Errorf("couldn't fetch current user from db: %w", err)
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), currUser.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows from current user: %w", err)
	}

	fmt.Printf("%s's feeds:\n", currUser.Name)
	fmt.Println("========")
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}
