package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	commandFunc, exists := c.registeredCommands[cmd.Name]
	if exists == false {
		return errors.New("command doesn't exist")
	}
	return commandFunc(s, cmd)
}

func NewCommands() commands {
	return commands{registeredCommands: make(map[string]func(*state, command) error)}
}
