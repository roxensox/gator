package state

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/roxensox/gator/internal/database"
	"github.com/roxensox/gator/internal/rss"
	"os"
	"time"
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

func HandlerAgg(s *State, cmd Command) error {
	// Handles the agg command

	// NOTE: Will be updated

	// Reads the feed into an RSSFeed object
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	// Returns any errors
	if err != nil {
		return err
	}

	// Prints the whole RSSFeed object
	fmt.Println(feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	// Handles the addfeed command

	if len(cmd.Args) < 2 {
		fmt.Println("Must supply a name for the feed and a target URL.")
		os.Exit(1)
	}

	// Gets the current time
	currTime := time.Now()

	// Queries the current user from the DB
	user, err := s.Conn.GetUser(context.Background(), *s.Cfg.CurrentUser)
	// Reports any errors and exits with error code
	if err != nil {
		fmt.Println("Failed to get user ID.")
		fmt.Printf("Error: \n\t%s\n", err)
		os.Exit(1)
	}

	// Creates the parameter object for CreateFeed
	inVal := database.CreateFeedParams{
		// Generates a uuid for the feed
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time:  currTime,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  currTime,
			Valid: true,
		},
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
		UserID: user.ID,
	}

	// Adds the new feed to the database
	newFeed, err := s.Conn.CreateFeed(context.Background(), inVal)

	// Reports any errors and exits with failure code
	if err != nil {
		fmt.Println("Failed to add feed.")
		fmt.Printf("Error: \n\t%s\n", err)
		os.Exit(1)
	}

	// Prints details of the new feed
	fmt.Printf(
		"Name: %s\nURL: %s\nFeed ID: %s\nUser ID: %s\nCreated at: %s\nUpdated at: %s\n",
		newFeed.Name,
		newFeed.Url,
		newFeed.ID,
		newFeed.UserID,
		newFeed.CreatedAt.Time,
		newFeed.UpdatedAt.Time,
	)

	return nil
}

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

	// Creates the user
	usr, err := s.Conn.CreateUser(context.Background(), newUser)

	// If an error results, returns it
	if err != nil {
		return err
	}

	// Sets the current user to the new user
	s.Cfg.SetUser(cmd.Args[0])

	// Prints a success message with details [DEBUG ONLY]
	//fmt.Printf("User: %s created\nUUID: %s\nCreated At: %s\nUpdated At: %s\n",
	//cmd.Args[0],
	//newUser.ID,
	//newUser.CreatedAt.Time,
	//newUser.UpdatedAt.Time,
	//)

	// Prints a regular success message
	fmt.Printf("Successfully registered %s\n", cmd.Args[0])

	// Returns a null error
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
