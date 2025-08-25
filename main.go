package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/groovypotato/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdList map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required")
	}
	username := cmd.args[0]
	if err := s.config.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("user has been set to %s", username)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	if value, ok := c.cmdList[cmd.name]; ok {
		return value(s, cmd)
	} else {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.cmdList == nil {
		c.cmdList = make(map[string]func(*state, command) error)
	}
	c.cmdList[name] = f
}

func main() {
	currConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = currConfig.SetUser("groovypotato")
	if err != nil {
		fmt.Println(err)
		return
	}
	newConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		jsonData, _ := json.MarshalIndent(newConfig, "", "  ")
		fmt.Println(string(jsonData))
	}
}
