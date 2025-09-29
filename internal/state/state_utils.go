package state

import (
	"context"
	"fmt"
	"github.com/roxensox/gator/internal/database"
	"strconv"
)

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	// The lengths you have to go to for conditional assignment
	limit := func(args []string) int {
		if len(args) > 0 {
			outVal, err := strconv.Atoi(cmd.Args[0])
			if err != nil {
				return 2
			}
			return outVal
		}
		return 2
	}(cmd.Args)

	fmt.Printf("Showing %d results:\n\n", limit)

	params := database.GetPostsForUserParams{
		UID:   user.ID,
		Limit: int32(limit),
	}

	results, err := s.Conn.GetPostsForUser(context.Background(), params)
	if err != nil {
		fmt.Printf("Could not get posts for %s.\n", user.Name)
		return err
	}

	for _, r := range results {
		fmt.Println(r.Title.String)
		fmt.Printf("%s\n\n", r.Url.String)
	}
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	// Runs a registered command

	// Gets the name of the command
	name := cmd.Name

	// Gets the corresponding function from the registry
	C, ok := c.Reg[name]

	// If there is no corresponding function, returns a custom error
	if !ok {
		return fmt.Errorf("Unknown command: %s\n", name)
	}

	// Calls the function and returns its results
	return C(s, cmd)
}

func (c *Commands) Register(name string, f func(s *State, cmd Command) error) {
	// Registers a command with a corresponding function

	// Checks if the command already has a registered function, returns early if so
	if _, ok := c.Reg[name]; ok {
		fmt.Println("Function already registered under this name.")
		return
	}

	// Adds the function to the registry
	c.Reg[name] = f
}
