package main

import (
	"fmt"
	"log"

	"github.com/NikolaTosic-sudo/gator/internal/config"
	"github.com/NikolaTosic-sudo/gator/internal/database"
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
