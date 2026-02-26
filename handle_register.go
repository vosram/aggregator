package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vosram/aggregator/internal/database"
)

func handleRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <username>\n", cmd.Name)
	}
	name := cmd.Args[0]
	now := time.Now().UTC()
	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	}
	user, err := s.db.CreateUser(context.Background(), args)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w\n", err)
	}
	fmt.Printf("Created: %s with ID %v, successfully!\n", user.Name, user.ID)
	err = s.conf.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set user: %w\n", err)
	}
	return nil
}
