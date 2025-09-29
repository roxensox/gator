package state

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/roxensox/gator/internal/database"
	"github.com/roxensox/gator/internal/rss"
	"log"
	"os"
	"strings"
	"time"
)

func HandlerAgg(s *State, cmd Command, user database.User) error {
	// Handles the agg command

	if len(cmd.Args) < 1 {
		fmt.Println("Must provide a time between requests interval")
		os.Exit(1)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		fmt.Printf("Unable to convert %s to a time duration.\n\nError: %s", cmd.Args[0], err)
		os.Exit(1)
	}

	ticker := time.NewTicker(timeBetweenReqs)
	log.Printf("Refreshing feeds every %v.", timeBetweenReqs)
	for ; ; <-ticker.C {
		ScrapeFeeds(s)
	}
}

func ScrapeFeeds(s *State) {
	// Scrapes RSS Feeds from their URLs

	// Gets the feed with the oldest last update from the db
	feed, err := s.Conn.GetFeedToCheck(context.Background())
	if err != nil {
		fmt.Printf("Unable to get feed\n")
		return
	}
	scrapeFeed(s.Conn, feed)
}

func scrapeFeed(conn *database.Queries, feed database.Feed) {

	// Builds params for MarkFeedChecked
	params := database.MarkFeedCheckedParams{
		LastChecked: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		ID: feed.ID,
	}
	err := conn.MarkFeedChecked(
		context.Background(),
		params,
	)

	if err != nil {
		fmt.Println("Unable to mark feed checked.")
		return
	}

	// Reads the feed into an RSSFeed object
	content, err := rss.FetchFeed(context.Background(), feed.Url)
	// Returns any errors
	if err != nil {
		fmt.Println("Failed to fetch content.")
		fmt.Printf("Error: %s\n", err)
		return
	}
	if content == nil {
		return
	}

	// Prints each item in the feed's item list
	for _, i := range content.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, i.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		params := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			Title: sql.NullString{
				String: i.Title,
				Valid:  true,
			},
			Url: sql.NullString{
				String: i.Link,
				Valid:  true,
			},
			Description: sql.NullString{
				String: i.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		}
		_, err = conn.CreatePost(
			context.Background(),
			params,
		)
		if err != nil {
			if !strings.Contains(fmt.Sprintf("%s", err), "duplicate key value") {
				log.Printf("%v", err)
			}
		}
	}
}
