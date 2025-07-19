package main

import (
	"log"
	"os"

	"github.com/NikolaTosic-sudo/gator/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()

	if err != nil {
		log.Fatalf("error reading config %v", err)
	}

	s := &State{
		Cfg: cfg,
	}

	cmds := commands{
		command: make(map[string]func(*State, cliCommand) error),
	}

	cmds.register("login", handleLogin)

	arguments := os.Args

	if len(arguments) < 2 {
		log.Fatalf("not enough arguments \n")
	}

	cmds.run(s, cliCommand{
		name:      arguments[1],
		arguments: arguments[2:],
		callback:  cmds.command[arguments[1]],
	})
}
