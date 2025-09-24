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
	usr, err := s.Conn.GetUser(context.Background(), cmd.Args[0])
	if err != nil || usr.Name != cmd.Args[0] {
		fmt.Printf("User %s not found.\n", cmd.Args[0])
		os.Exit(1)
	}
	err = s.Cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set to %s.\n", cmd.Args[0])
	return nil
}

func HandlerReset(s *State, _ Command) error {
	err := s.Conn.ResetUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
	return nil
}

func HandlerGetUsers(s *State, _ Command) error {
	usrs, err := s.Conn.GetUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}
	for _, v := range usrs {
		if v == *s.Cfg.CurrentUser {
			fmt.Printf(" * %s (current)\n", v)
		} else {
			fmt.Printf(" * %s\n", v)
		}
	}
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Command has no arguments")
	}
	usr, _ := s.Conn.GetUser(context.Background(), cmd.Args[0])
	if usr.Name == cmd.Args[0] {
		fmt.Printf("User %s already registered.\n", cmd.Args[0])
		os.Exit(1)
	}
	currTime := time.Now()
	newUser := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time:  currTime,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  currTime,
			Valid: true,
		},
		Name: cmd.Args[0],
	}
	usr, err := s.Conn.CreateUser(context.Background(), newUser)
	if err != nil {
		return err
	}
	s.Cfg.SetUser(cmd.Args[0])
	fmt.Printf("User: %s created\nUUID: %s\nCreated At: %s\nUpdated At: %s\n",
		cmd.Args[0],
		newUser.ID,
		newUser.CreatedAt.Time,
		newUser.UpdatedAt.Time,
	)
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
