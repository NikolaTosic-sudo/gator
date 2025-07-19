package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/NikolaTosic-sudo/gator/internal/config"
	"github.com/NikolaTosic-sudo/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()

	if err != nil {
		log.Fatalf("error reading config %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)

	if err != nil {
		log.Fatalf("error opening the database %v", err)
	}

	dbQueries := database.New(db)

	s := &State{
		db:  dbQueries,
		Cfg: &cfg,
	}

	cmds := commands{
		command: make(map[string]func(*State, cliCommand) error),
	}

	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)

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
