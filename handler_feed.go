package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/enderbd/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)


func handleAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Missing argument, usage: %s <time_between_reqs>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error trying to create the timeBetweenRequest duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}


func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("missing argument(s), usage: %s <feed_name> <feed_url>", cmd.name)
	}

	
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),	
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't not add the feed: %w", err)
	}

	fmt.Println("New feed created ")
	printFeed(feed)

	_ , err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Could not create an entry into feed_follow")
	}

	fmt.Println("Feed added to the feed follows as well !")

	return nil
}

func handleFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error trying to get all the feeds: %w", err)	
	}

	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.Name)
		fmt.Printf("* %s\n", feed.Url)
		userId := feed.UserID
		creator, err := s.db.GetUserWithId(context.Background(), userId)
		if err != nil {
			return fmt.Errorf("Could not find the feed creator user %w", err)
		}
		fmt.Printf("* %s\n", creator.Name)

	}

	return nil
}


func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())	
	if err != nil{
		log.Printf("Error trying to fetch next feed: %v\n", err)
		return
	}
	

	rssFeeds, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Error trying to fetch the rss feed: %v\n", err)
		return
	}

	for _, item := range rssFeeds.Channel.Item {
		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		validTime := true
		if err != nil {
			t = time.Time{}
			validTime = false
		}

		pubDate := sql.NullTime {
			Time: t,
			Valid: validTime,
		}

		title := sql.NullString{
			String: item.Title,
			Valid: item.Title != "",
		}

		description := sql.NullString{
			String: item.Description,
			Valid: item.Description != "",
		}

		_ , err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: title,
			Url: item.Link,
			Description: description,
			PublishedAt: pubDate,
			FeedID: feed.ID,
		})	
		if err != nil {
			 // Attempt to cast the error to a pq.Error type
		    	if pqErr, ok := err.(*pq.Error); ok {
				// Check if the SQLSTATE code indicates a unique violation
				if pqErr.Code == "23505" {
			    	// This is a unique constraint violation, ignore it as per assignment
			    		continue // Skip to the next item in the loop
				}
		    	}
			log.Printf("Error adding a post %v\n", err)
		}
	}

	if err = s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		log.Printf("Error trying to mark feed fetched: %v\n", err)
		return
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeeds.Channel.Item))

}
