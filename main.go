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
	cfg := config.NewConfig()
	cfg.Read()
	s := state.State{
		Cfg: &cfg,
	}
	db, err := sql.Open("postgres", *s.Cfg.DB_URL)
	if err != nil {
		fmt.Println("Unable to access database")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s.Conn = dbQueries

	cmds := state.Commands{
		Reg: make(map[string]func(s *state.State, c state.Command) error),
	}

	cmds.Register("login", state.HandlerLogin)
	cmds.Register("register", state.HandlerRegister)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Must specify a command.")
		os.Exit(1)
	}

	if args[1] == "login" && len(args) < 3 {
		fmt.Println("Must provide a username to log in.")
		os.Exit(1)
	}

	cmd := state.Command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.Run(&s, cmd)

	if err != nil {
		fmt.Print(err)
	}
	os.Exit(0)
}
