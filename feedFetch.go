package main

import(
	"fmt"
	"errors"
	"context"
	"io"
	"time"
	"html"
	"net/http"
	"encoding/xml"
	"database/sql"

	"github.com/google/uuid"

	"github.com/F0hor/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func scrapeFeeds(s *state) (error) {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Failed to retrieve data from database:\n%v\n", err)
	}

	err = s.db.MarkFeedFetched(
		ctx,
		database.MarkFeedFetchedParams{
			ID: feed.ID,
			UpdatedAt: time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("Failed to update fetch time:\n%v\n", err)
	}

	rss, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("Something went wrong with feed request:\n%v\n", err)
	}

	fmt.Printf("RSS Feed from %v (%v):\n\n", rss.Channel.Title, rss.Channel.Link)
	for _, i := range rss.Channel.Items {
		fmt.Printf(" - %v (%v at %v)\n", i.Title, i.PubDate, i.Link)

		var pubAt sql.NullTime
		pub, err := time.Parse(time.RFC1123Z, i.PubDate)
		if err != nil {
			pubAt.Valid = false
		} else {
			pubAt.Time = pub
			pubAt.Valid = true
		}

		var desc sql.NullString
		if len(i.Description) == 0 {
			desc.Valid = false
		} else {
			desc.String = i.Description
		}

		_, err = s.db.CreatePost(
			ctx,
			database.CreatePostParams{
				ID: uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Title: i.Title,
				Url: i.Link,
				Description: desc,
				PublishedAt: pubAt,
				FeedID: feed.ID,
			},
		)
		if err != nil {
			fmt.Errorf("Failed to insert post into database:\n%v\n", err)
		}
	}

	return nil
}

func fetchFeed(ctx context.Context, url string) (*RSSFeed, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Request to %s failed: %v\n", url, err)
	}

	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request to %s failed: %v\n", url, err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body:\n%v\n", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, errors.New("Failed to decode the data")
	}

	unescapeAll(&feed)
	
	return &feed, nil
}

func unescapeAll(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, i := range feed.Channel.Items {
		i.Title = html.UnescapeString(i.Title)
		i.Description = html.UnescapeString(i.Description)
	}
}

