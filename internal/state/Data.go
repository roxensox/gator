package state

import (
	"context"
	"os"
)

func HandlerReset(s *State, _ Command) error {
	// Handles the reset command

	// Wipes the users table
	err := s.Conn.ResetUsers(context.Background())

	// If there is an error wiping the table, exits with code 1
	if err != nil {
		os.Exit(1)
	}

	// Exits with code 0
	os.Exit(0)

	// Returns null error so this will run
	return nil
}
