package state

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/roxensox/gator/internal/database"
	"os"
	"time"
)

func HandlerFollow(s *State, cmd Command, user database.User) error {
	// Handles the follow command

	// Gets the current time
	currTime := time.Now().UTC()

	// Queries the feed from the database
	feed, err := s.Conn.GetFeed(context.Background(), cmd.Args[0])
	// Reports any errors and exits with code 1
	if err != nil {
		fmt.Printf("Unable to find feed for %s\nError:\n\t%s\n",
			cmd.Args[0],
			err,
		)
	}

	// Creates a CreateFeedFollow paremeter object
	inVal := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currTime,
		UpdatedAt: currTime,
		UID:       user.ID,
		FID:       feed.ID,
	}

	// Creates the feed follow
	feed_follow, err := s.Conn.CreateFeedFollow(context.Background(), inVal)
	// Reports any errors and exits with code 1
	if err != nil {
		fmt.Printf("Unable to follow feed.\nError:\n\t%s\n", err)
		os.Exit(1)
	}

	// Prints the follow details
	fmt.Printf(
		"Name: %s\n\tURL: %s\n\tUser: %s\n",
		feed_follow.FeedName,
		feed_follow.FeedUrl,
		feed_follow.UserName,
	)
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	// Handles the unfollow command
	if len(cmd.Args) < 1 {
		fmt.Println("Must provide a url to unfollow.")
		os.Exit(1)
	}

	feed, err := s.Conn.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Printf("Unable to find feed for %s\nError:\n\t%s\n",
			cmd.Args[0],
			err,
		)
	}
	params := database.UnfollowFeedParams{
		UID: user.ID,
		FID: feed.ID,
	}
	err = s.Conn.UnfollowFeed(context.Background(), params)
	if err != nil {
		fmt.Printf("Unable to unfollow feed %s\nError:\n\t%s",
			feed.Name,
			err,
		)
		os.Exit(1)
	}
	return nil
}

func HandlerFollowing(s *State, _ Command, user database.User) error {
	// Handles the following command

	// Queries user follows from the database
	follows, err := s.Conn.GetFeedFollowsForUser(context.Background(), user.ID)
	// Reports any errors and exits with code 1
	if err != nil {
		fmt.Printf("Unable to get user follows.\nError:\n\t%s\n", err)
		os.Exit(1)
	}

	// Iterates through all feed follows
	for i, f := range follows {
		// Prints the details
		fmt.Printf(
			"Feed %d:\n\tName: %s\n\tURL: %s\n\tUser: %s\n",
			i+1,
			f.FeedName,
			f.FeedUrl,
			f.UserName,
		)
	}
	return nil
}
