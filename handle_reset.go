package main

import (
	"context"
	"fmt"
)

func handleReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not delete all users: %w", err)
	}
	fmt.Println("All users deleted successfully")
	return nil
}
