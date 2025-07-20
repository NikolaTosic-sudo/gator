package main

import (
	"context"
	"fmt"
	"html"
	"time"

	"github.com/NikolaTosic-sudo/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *State, cmd cliCommand, user database.User) error) func(*State, cliCommand) error {

	// var emptyUser database.User

	return handleLogin
}

func scrapeFeeds(s *State) error {
	client := &Client{}
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return err
	}

	rssFeed, err := client.fetchFeed(context.Background(), nextFeed.Url)

	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        nextFeed.ID,
	})

	if err != nil {
		return err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	fmt.Println(rssFeed.Channel.Title)

	for i, post := range rssFeed.Channel.Item {
		post.Title = html.UnescapeString(post.Title)
		post.Description = html.UnescapeString(post.Description)
		rssFeed.Channel.Item[i] = post
		fmt.Println(rssFeed.Channel.Item[i].Title)
	}

	return nil
}
