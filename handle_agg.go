package main

import (
	"fmt"
	"time"
)

func handleAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg < Duration: 1s | 1m | 1h >\nEx. agg 1m")
	}
	timeString := cmd.Args[0]
	timeBetweenRequest, err := time.ParseDuration(timeString)
	if err != nil {
		return fmt.Errorf("couldn't parse duration time: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Println(err)
		}
	}
}
