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
	cmds.register("reset", handleReset)
	cmds.register("users", handleGetAllUsers)
	cmds.register("agg", handleFetch)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleGetAllFeeds)
	cmds.register("follow", middlewareLoggedIn(handleCreateFeedFollow))
	cmds.register("following", middlewareLoggedIn(handleGetFeedFollowsForUser))
	cmds.register("unfollow", middlewareLoggedIn(handleRemoveFeedFollow))
	cmds.register("browse", middlewareLoggedIn(handleBrowse))

	arguments := os.Args

	if len(arguments) < 2 {
		log.Fatalf("not enough arguments \n")
	}

	err = cmds.run(s, cliCommand{
		name:      arguments[1],
		arguments: arguments[2:],
		callback:  cmds.command[arguments[1]],
	})

	if err != nil {
		log.Fatal(err)
	}
}
