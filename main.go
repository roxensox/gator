package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/roxensox/gator/internal/config"
	"github.com/roxensox/gator/internal/database"
	"github.com/roxensox/gator/internal/state"
	"os"
)

func main() {
	// Initializes the configuration object
	cfg := config.NewConfig()

	// Reads the config file into the object
	cfg.Read()

	// Adds config to a new state object
	s := state.State{
		Cfg: &cfg,
	}

	// Creates a connection to the database
	db, err := sql.Open("postgres", *s.Cfg.DB_URL)

	// If there is an error, exits with code 1
	if err != nil {
		fmt.Println("Unable to access database")
		os.Exit(1)
	}

	// Creates a query engine
	dbQueries := database.New(db)

	// Adds the query engine to the state
	s.Conn = dbQueries

	// Initializes the command registry
	cmds := state.Commands{
		Reg: make(map[string]func(s *state.State, c state.Command) error),
	}

	// Registers command:handler pairs
	cmds.Register("login", state.HandlerLogin)
	cmds.Register("register", state.HandlerRegister)
	cmds.Register("reset", state.HandlerReset)
	cmds.Register("users", state.HandlerGetUsers)
	cmds.Register("agg", state.MiddlewareLoggedIn(state.HandlerAgg))
	cmds.Register("addfeed", state.MiddlewareLoggedIn(state.HandlerAddFeed))
	cmds.Register("feeds", state.MiddlewareLoggedIn(state.HandlerFeeds))
	cmds.Register("follow", state.MiddlewareLoggedIn(state.HandlerFollow))
	cmds.Register("following", state.MiddlewareLoggedIn(state.HandlerFollowing))
	cmds.Register("unfollow", state.MiddlewareLoggedIn(state.HandlerUnfollow))

	// Reads input into args
	args := os.Args

	// Exits with code 1 if user supplies no arguments
	if len(args) < 2 {
		fmt.Println("Must specify a command.")
		os.Exit(1)
	}

	// Splits the arguments into a command
	cmd := state.Command{
		Name: args[1],
		Args: args[2:],
	}

	// Runs the command
	err = cmds.Run(&s, cmd)

	// Prints error and exits with code 1 if there is an error
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	// Exits with code 0
	os.Exit(0)
}
