package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/enderbd/gator/internal/config"
	"github.com/enderbd/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file %v", err)
	}
	
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	toolState := &state{
		db: dbQueries,
		cfg: &cfg,
	}
	
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleUsers)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(toolState, command{name: cmdName, args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}
