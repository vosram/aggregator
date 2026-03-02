package main

import (
	"context"
	"fmt"

	"github.com/vosram/aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.conf.CurrentUser)
		if err != nil {
			return fmt.Errorf("couldn't get logged in user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
