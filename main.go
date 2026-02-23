package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vosram/aggregator/internal/config"
)

type state struct {
	conf *config.Config
}

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Println("Error:", err)
	}
	state := state{conf: &config}
	activeCommands := NewCommands()
	activeCommands.register("login", handleLogin)
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

	err = activeCommands.run(&state, cliCommand)
	if err != nil {
		log.Fatal("Error: ", err)
	}

}
