package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NikolaTosic-sudo/gator/internal/config"
	"github.com/NikolaTosic-sudo/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	db  *database.Queries
	Cfg *config.Config
}

type cliCommand struct {
	name      string
	arguments []string
	callback  func(state *State, cmd cliCommand) error
}

type commands struct {
	command map[string]func(*State, cliCommand) error
}

func handleLogin(state *State, cmd cliCommand) error {
	if len(cmd.arguments) == 0 {
		log.Fatal("please enter your username")
		return fmt.Errorf("please enter your username")
	}

	user, err := state.db.GetUser(context.Background(), cmd.arguments[0])

	if err != nil {
		log.Fatal("user with this username doesn't exist")
		return err
	}

	err = state.Cfg.SetUser(user.Name)

	if err != nil {
		log.Fatal("invalid username")
		return fmt.Errorf("invalid username")
	}

	return nil
}

func handleRegister(state *State, cmd cliCommand) error {

	if len(cmd.arguments) == 0 {
		log.Fatal("please enter a name")
		return fmt.Errorf("please enter a name")
	}

	user, _ := state.db.GetUser(context.Background(), cmd.arguments[0])

	var emptyUser database.User

	if user != emptyUser {
		log.Fatal("User with that name already exists")
		return nil
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
		fmt.Printf("there was an error with creating the user %v \n", err)
		return fmt.Errorf("err")
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
		log.Fatal("couldn't get users")
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

func (c *commands) run(s *State, cmd cliCommand) error {

	if s == nil {
		log.Fatal("state failed, please restart the application")
		return fmt.Errorf("state failed, please restart the application")
	}

	command := c.command[cmd.name]

	command(s, cmd)

	return nil
}

func (c *commands) register(name string, function func(*State, cliCommand) error) error {

	if name == "" {
		log.Fatal("please enter command's name")
		return fmt.Errorf("please enter command's name")
	}

	c.command[name] = function

	return nil
}
