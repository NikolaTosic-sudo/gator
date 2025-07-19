package main

import (
	"fmt"

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
		return fmt.Errorf("state failed, please restart the application")
	}

	command, ok := c.command[cmd.name]

	if !ok {
		return fmt.Errorf("unknown command")
	}

	err := command(s, cmd)

	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, function func(*State, cliCommand) error) error {

	if name == "" {
		return fmt.Errorf("please enter command's name")
	}

	c.command[name] = function

	return nil
}
