package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/groovypotato/gator/internal/config"
	"github.com/groovypotato/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
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
	u, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("user not found:", err)
		os.Exit(1)
	}
	if err := s.config.SetUser(u.Name); err != nil {
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

func handlerRegister(s *state, cmd command) error {
    if len(cmd.args) == 0 {
        return errors.New("username is required")
    }
    username := cmd.args[0]

    p := database.CreateUserParams{
        ID:        uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name:      username,
    }

    u, err := s.db.CreateUser(context.Background(), p)
    if err != nil {
        fmt.Println("error creating user:", err)
        os.Exit(1)
    }

    if err := s.config.SetUser(username); err != nil {
        return err
    }

    fmt.Printf("user created: %+v\n", u)
    return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
        return errors.New("too many arguments")
    }

	err := s.db.DeleteUsers(context.Background())
    if err != nil {
        fmt.Println("error resetting users:", err)
        os.Exit(1)
    }
	fmt.Println("users table has been reset")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
        return errors.New("too many arguments")
    }
	u,err := s.db.GetUsers(context.Background())
    if err != nil {
        fmt.Println("error resetting users:", err)
        os.Exit(1)
    }
	for _,user := range u {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		}else {
			fmt.Printf("* %s\n",user.Name)
		}
	}
	return nil
}

func buildRequest(ctx context.Context, feedURL string) (*http.Request, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", "gator")
    return req, nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}
	newCommand := command{
		name: os.Args[1],
		args: []string{},
	}
	if len(os.Args) > 2 {
		newCommand.args = os.Args[2:]
	}
	var newState state
	currConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	newState.config = &currConfig
	newCommands := commands{
		cmdList: make(map[string]func(*state, command) error),
	}
	newCommands.register("login", handlerLogin)
	newCommands.register("register", handlerRegister)
	newCommands.register("reset", handlerReset)
	newCommands.register("users", handlerGetUsers)

	dbURL := newState.config.DBURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)
	newState.db = dbQueries
	newCommands.run(&newState, newCommand)
	newConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		jsonData, _ := json.MarshalIndent(newConfig, "", "  ")
		fmt.Println(string(jsonData))
	}
}
