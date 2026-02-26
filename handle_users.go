package main

import (
	"context"
	"fmt"
)

func handleGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	currentUser := s.conf.CurrentUser
	if err != nil {
		return fmt.Errorf("couldn't fetch all users: %w", err)
	}
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
