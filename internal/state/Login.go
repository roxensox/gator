package state

import (
	"context"
	"fmt"
	"os"
)

func HandlerLogin(s *State, cmd Command) error {
	// Handles the login command

	// If the user didn't provide a username to log in, returns an error
	if len(cmd.Args) == 0 {
		fmt.Printf("Must provide a username to log in.\n")
		os.Exit(1)
	}

	// Queries the user from the database
	usr, err := s.Conn.GetUser(context.Background(), cmd.Args[0])

	// If the database doesn't return anything, exits with code 1
	if err != nil || usr.Name != cmd.Args[0] {
		fmt.Printf("User %s not found.\n", cmd.Args[0])
		os.Exit(1)
	}

	// Sets the user to the provided username
	err = s.Cfg.SetUser(cmd.Args[0])

	// If there is an error setting the user, returns it
	if err != nil {
		return err
	}

	// Prints success message
	fmt.Printf("User set to %s.\n", cmd.Args[0])

	// Returns a null error
	return nil
}
