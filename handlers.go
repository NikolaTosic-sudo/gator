package main

import (
	"context"
	"fmt"
	"html"
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
	client := &Client{}
	rssFeed, err := client.fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if err != nil {
		return err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i, post := range rssFeed.Channel.Item {
		post.Title = html.UnescapeString(post.Title)
		post.Description = html.UnescapeString(post.Description)
		rssFeed.Channel.Item[i] = post
	}

	fmt.Print(rssFeed)

	return nil
}
