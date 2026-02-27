package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/vosram/aggregator/internal/config"
	"github.com/vosram/aggregator/internal/database"
)

type state struct {
	db   *database.Queries
	conf *config.Config
}

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Println("Error:", err)
	}
	activeCommands := NewCommands()
	activeCommands.register("login", handleLogin)
	activeCommands.register("register", handleRegister)
	activeCommands.register("reset", handleReset)
	activeCommands.register("users", handleGetUsers)
	activeCommands.register("agg", handleAgg)

	db, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	state := &state{conf: &config, db: dbQueries}

	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("Not enough command line arguments:", args)
	}
	cliCommand := command{
		Args: make([]string, 0),
	}
	for i, elem := range args {
		if i == 0 {
			cliCommand.Name = elem
		} else {
			cliCommand.Args = append(cliCommand.Args, elem)
		}
	}

	err = activeCommands.run(state, cliCommand)
	if err != nil {
		log.Fatal("Error: ", err)
	}

}
