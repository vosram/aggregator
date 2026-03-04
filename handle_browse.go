package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/vosram/aggregator/internal/database"
)

func handleBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	var err error
	if len(cmd.Args) == 1 {
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("couldn't convert limit arg to an int: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get post for user: %w", err)
	}
	if len(posts) == 0 {
		fmt.Printf("No post saved for %s", user.Name)
		return nil
	}
	fmt.Printf("%s's posts:\n=======\n", user.Name)
	for _, post := range posts {
		fmt.Printf("Feed: %s\n", post.FeedName)
		fmt.Println(post.Title)
		fmt.Println(post.PublishedAt.Time.Format(time.RFC1123Z))
		fmt.Println("    ", post.Description)
		fmt.Println("Link: ", post.Url)
		fmt.Println("")
	}
	return nil
}
