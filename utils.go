package main

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	"github.com/NikolaTosic-sudo/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *State, cmd cliCommand, user database.User) error) func(*State, cliCommand) error {

	return func(s *State, cmd cliCommand) error {

		user, err := s.db.GetUser(context.Background(), s.Cfg.CurrentUserName)

		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
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
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, post.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		if err != nil {
			return err
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       post.Title,
			Url:         post.Title,
			Description: sql.NullString{String: post.Description, Valid: post.Description != ""},
			PublishedAt: publishedAt,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	log.Printf("Feed %s collected, %v posts found", rssFeed.Channel.Title, len(rssFeed.Channel.Item))

	return nil
}
