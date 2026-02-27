package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
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
