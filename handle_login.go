package main

import (
	"context"
	"fmt"
	"log"
)

func handleLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <username>\n", cmd.Name)
	}
	name := cmd.Args[0]

	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		log.Fatal("user doesn't exist in database")
	}

	err = s.conf.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w\n", err)
	}

	fmt.Printf("User %s has been set successfully!\n", name)
	return nil
}
