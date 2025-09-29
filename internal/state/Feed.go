package state

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/roxensox/gator/internal/database"
	"os"
	"time"
)

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	// Handles the addfeed command

	// Validates length of input
	if len(cmd.Args) < 2 {
		fmt.Println("Must supply a name for the feed and a target URL.")
		os.Exit(1)
	}

	// Gets the current time
	currTime := time.Now().UTC()

	// Creates the parameter object for CreateFeed
	inVal := database.CreateFeedParams{
		// Generates a uuid for the feed
		ID:        uuid.New(),
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	// Adds the new feed to the database
	newFeed, err := s.Conn.CreateFeed(context.Background(), inVal)
	// Reports any errors and exits with failure code
	if err != nil {
		fmt.Println("Failed to add feed.")
		fmt.Printf("Error: \n\t%s\n", err)
		os.Exit(1)
	}

	// Creates a CreateFeedFollow paremeter object
	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UID:       user.ID,
		FID:       newFeed.ID,
		CreatedAt: currTime,
		UpdatedAt: currTime,
	}

	// Follows the feed for the currently signed in user
	_, err = s.Conn.CreateFeedFollow(context.Background(), feedFollow)
	// Reports any errors and exits with code 1
	if err != nil {
		fmt.Println("Failed to follow feed.")
		fmt.Printf("Error: \n\t%s\n", err)
		os.Exit(1)
	}

	// Prints details of the new feed
	fmt.Printf(
		"Name: %s\nURL: %s\n",
		newFeed.Name,
		newFeed.Url,
	)

	return nil
}

func HandlerFeeds(s *State, _ Command, user database.User) error {
	// Handles the feeds command

	// Queries the feeds from the database
	results, err := s.Conn.GetFeeds(context.Background())
	// Reports any errors and exits with error code
	if err != nil {
		fmt.Println("Failed to get feed information.")
		fmt.Printf("Error: \n\t%s\n", err)
		os.Exit(1)
	}

	// Creates an empty slice of errors
	errors := []error{}

	// Loops through results
	for i, feed := range results {
		// Queries the user from the feed by ID
		usr, err := s.Conn.GetUserFromId(context.Background(), feed.UserID)
		// If there's a problem getting the user, adds the error to the errors slice
		if err != nil {
			errors = append(errors, err)
			continue
		}
		// Prints the feed information
		fmt.Printf("Feed %d:\n\tName: %s\n\tURL: %s\n\tUser: %s\n",
			i+1,
			feed.Name,
			feed.Url,
			usr.Name,
		)
	}

	// Loops through and prints errors
	if len(errors) > 0 {
		fmt.Printf("Errors:\n\t")
		for _, e := range errors {
			fmt.Printf("%s\n\t", e)
		}
		// Exits with error code
		os.Exit(1)
	}

	return nil
}
