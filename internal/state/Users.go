package state

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/roxensox/gator/internal/database"
	"os"
	"time"
)

func HandlerGetUsers(s *State, _ Command) error {
	// Handles the users command

	// Queries all usernames from the database
	usrs, err := s.Conn.GetUsers(context.Background())

	// If there is an error from the query, exits with code 1
	if err != nil {
		os.Exit(1)
	}

	// Iterates through all usernames
	for _, v := range usrs {
		// If the user being read is logged in, prints special message
		// Otherwise, prints the username
		if v == *s.Cfg.CurrentUser {
			fmt.Printf(" * %s (current)\n", v)
		} else {
			fmt.Printf(" * %s\n", v)
		}
	}

	// Returns a null error
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	// Handles the register command

	// If the user didn't provide a username, returns an error
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Command has no arguments")
	}

	// Queries the database to check if the user has already been added
	usr, _ := s.Conn.GetUser(context.Background(), cmd.Args[0])

	// If a user is returned and matches the input, exits with code 1
	if usr.Name == cmd.Args[0] {
		fmt.Printf("User %s already registered.\n", cmd.Args[0])
		os.Exit(1)
	}

	// Gets the current time
	currTime := time.Now()

	// Builds a CreateUserParams object for the input username
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Name:      cmd.Args[0],
	}

	// Creates the user
	usr, err := s.Conn.CreateUser(context.Background(), newUser)

	// If an error results, returns it
	if err != nil {
		return err
	}

	// Sets the current user to the new user
	s.Cfg.SetUser(cmd.Args[0])

	// Prints a regular success message
	fmt.Printf("Successfully registered %s\n", cmd.Args[0])

	// Returns a null error
	return nil
}
