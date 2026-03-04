package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vosram/aggregator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// prepare request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	// do request
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request to %s: %w", feedURL, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is not 200; status code %d", res.StatusCode)
	}

	// parse response
	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var xmlData RSSFeed
	if err = xml.Unmarshal(dat, &xmlData); err != nil {
		return nil, fmt.Errorf("error unmarshalling xml data: %w", err)
	}
	xmlData.Channel.Title = html.UnescapeString(xmlData.Channel.Title)
	xmlData.Channel.Description = html.UnescapeString(xmlData.Channel.Description)
	for i, item := range xmlData.Channel.Item {
		xmlData.Channel.Item[i].Title = html.UnescapeString(item.Title)
		xmlData.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &xmlData, nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		ID:            feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	fmt.Printf("\n%s's Posts\n", feedData.Channel.Title)
	fmt.Println("=======")
	for _, feedItem := range feedData.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, feedItem.PubDate)
		if err != nil {
			log.Printf("Couldn't parse the pubDate for %s's post: %s\n", feedData.Channel.Title, feedItem.Title)
			continue
		}

		existingPost, err := s.db.GetPostByUrl(context.Background(), feedItem.Link)
		if err == nil {
			// post already exist
			log.Printf("post %s already exist, skipping...\n", existingPost.Url)
			continue
		}

		now := time.Now().UTC()
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       feedItem.Title,
			Url:         feedItem.Link,
			Description: feedItem.Description,
			PublishedAt: publishedAt.UTC(),
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("couldn't create post or post already exist: %s from %s\n", feedItem.Title, feedData.Channel.Title)
			continue
		}
		fmt.Printf("Created Post for: %s\n", post.Title)
	}
	return nil
}
