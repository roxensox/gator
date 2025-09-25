package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create a new get request with the input URL and context
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	// Return any errors
	if err != nil {
		return nil, err
	}

	// Initialize an http client
	client := http.Client{}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Set the response's body to close when this function ends
	defer resp.Body.Close()

	// Return in case of bad status code
	if resp.StatusCode > 299 {
		return nil, err
	}

	// Read the response body
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Initialize an RSSFeed object
	out := RSSFeed{}

	// Unmarshal the data to the RSSFeed
	xml.Unmarshal(dat, &out)

	// Decodes escaped HTML in necessary fields
	html.UnescapeString(out.Channel.Title)
	html.UnescapeString(out.Channel.Description)
	for _, item := range out.Channel.Item {
		html.UnescapeString(item.Title)
		html.UnescapeString(item.Description)
	}

	// Return a pointer to the RSSFeed
	return &out, nil
}
