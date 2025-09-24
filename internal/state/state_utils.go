package state

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/roxensox/gator/internal/database"
	"os"
	"time"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Command has no arguments")
	}
	err := s.Cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set to %s.\n", cmd.Args[0])
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Command has no arguments")
	}
	usr, err := s.Conn.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	if usr.Name == cmd.Args[0] {
		fmt.Printf("User %s already registered.\n", cmd.Args[0])
		os.Exit(1)
	}
	newUser := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Name: cmd.Args[0],
	}
	usr, err = s.Conn.CreateUser(context.Background(), newUser)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	name := cmd.Name
	C, ok := c.Reg[name]
	if !ok {
		return fmt.Errorf("Unknown command: %s\n", name)
	}
	return C(s, cmd)
}

func (c *Commands) Register(name string, f func(s *State, cmd Command) error) {
	if _, ok := c.Reg[name]; ok {
		fmt.Println("Function already registered under this name.")
		return
	}
	c.Reg[name] = f
}
