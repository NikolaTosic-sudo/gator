package main

import (
	"context"
	"fmt"
	"time"

	"github.com/NikolaTosic-sudo/gator/internal/database"
	"github.com/google/uuid"
)

func handleLogin(state *State, cmd cliCommand) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("please enter your username")
	}

	user, err := state.db.GetUser(context.Background(), cmd.arguments[0])

	if err != nil {
		return fmt.Errorf("user with the username: %v doesn't exist", cmd.arguments[0])
	}

	if user.Name == state.Cfg.CurrentUserName {
		fmt.Printf("%v is already logged in", user.Name)
		return nil
	}

	err = state.Cfg.SetUser(user.Name)

	if err != nil {
		return fmt.Errorf("invalid username")
	}

	fmt.Printf("Login successful. \nWelcome %v", user.Name)

	return nil
}

func handleRegister(state *State, cmd cliCommand) error {

	if len(cmd.arguments) == 0 {
		return fmt.Errorf("please enter a name")
	}

	user, _ := state.db.GetUser(context.Background(), cmd.arguments[0])

	var emptyUser database.User

	if user != emptyUser {
		return fmt.Errorf("user with that name already exists")
	}

	createdUser, err := state.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.arguments[0],
		},
	)

	if err != nil {
		return fmt.Errorf("there was an error with creating the user %v", err)
	}

	err = state.Cfg.SetUser(createdUser.Name)

	if err != nil {
		return err
	}

	fmt.Print("The user was successfuly created\n")

	fmt.Print(createdUser)

	return nil
}

func handleReset(state *State, cmd cliCommand) error {
	err := state.db.ResetDB(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func handleGetAllUsers(state *State, cmd cliCommand) error {

	users, err := state.db.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == state.Cfg.CurrentUserName {
			fmt.Println(user.Name, "(current)")
		} else {
			fmt.Println(user.Name)
		}
	}

	return nil
}

func handleFetch(state *State, cmd cliCommand) error {

	if len(cmd.arguments) == 0 {
		return fmt.Errorf("please insert a desired time duration between feeds")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])

	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		err := scrapeFeeds(state)

		if err != nil {
			return err
		}
	}
}

func handleAddFeed(state *State, cmd cliCommand) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("please enter Name and URL")
	}

	user, err := state.db.GetUser(context.Background(), state.Cfg.CurrentUserName)

	if err != nil {
		return err
	}

	createFeed, err := state.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	})

	if err != nil {
		return err
	}

	_, err = state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    createFeed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Print(createFeed)

	return nil
}

func handleGetAllFeeds(state *State, cmd cliCommand) error {
	feeds, err := state.db.GetAllFeeds(context.Background())

	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := state.db.GetUserById(context.Background(), feed.UserID)

		if err != nil {
			return err
		}

		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(user.Name)
	}

	return nil
}

func handleCreateFeedFollow(state *State, cmd cliCommand) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("pleaes enter a valid url")
	}

	feed, err := state.db.GetFeedByUrl(context.Background(), cmd.arguments[0])

	if err != nil {
		return err
	}

	user, err := state.db.GetUser(context.Background(), state.Cfg.CurrentUserName)

	if err != nil {
		return err
	}

	followFeed, err := state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Println(followFeed.UserName)
	fmt.Println(followFeed.FeedName)

	return nil
}

func handleGetFeedFollowsForUser(state *State, cmd cliCommand) error {
	user, err := state.db.GetUser(context.Background(), state.Cfg.CurrentUserName)

	if err != nil {
		return err
	}

	followFeed, err := state.db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return err
	}

	for _, feed := range followFeed {
		fmt.Println(feed.FeedName)
	}

	fmt.Println(user.Name)

	return nil
}

func handleRemoveFeedFollow(state *State, cmd cliCommand) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("please enter feed's URL you want to unfollow")
	}

	user, err := state.db.GetUser(context.Background(), state.Cfg.CurrentUserName)

	if err != nil {
		return err
	}

	feed, err := state.db.GetFeedByUrl(context.Background(), cmd.arguments[0])

	if err != nil {
		return err
	}

	deletedFeed, err := state.db.RemoveFeedFollow(context.Background(), database.RemoveFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Print(deletedFeed)

	return nil
}
