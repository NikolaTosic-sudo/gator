package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	httpClient http.Client
}

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
}

func (c *Client) fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {

	var RSSFeed RSSFeed

	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)

	if err != nil {
		return &RSSFeed, err
	}

	req.Header.Set("User-Agent", "gator")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return &RSSFeed, err
	}

	if res.StatusCode > 299 {
		return &RSSFeed, fmt.Errorf("there was an issue with the request")
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return &RSSFeed, err
	}

	if err := xml.Unmarshal(data, &RSSFeed); err != nil {
		return &RSSFeed, err
	}

	return &RSSFeed, nil
}
