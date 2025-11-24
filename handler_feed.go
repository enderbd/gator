package main

import (
	"context"
	"fmt"
	"time"

	"github.com/enderbd/gator/internal/database"
	"github.com/google/uuid"
)


func handleAgg(s *state, _ command) error {
	feedUrl := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("error trying to get the rss feed: %w", err)
	}

	fmt.Printf("%+v\n", *feed)
	return nil
}


func handleAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("missing argument(s), usage: %s <feed_name> <feed_url>", cmd.name)
	}

	currentUser := s.cfg.CurrentUserName
	userInfo, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("Could not find current user %v in the database: %+w ", currentUser, err)
	}
	userId := userInfo.ID
	
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),	
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: userId,
	})

	if err != nil {
		return fmt.Errorf("couldn't not add the feed: %w", err)
	}

	fmt.Println("New feed created ")
	printFeed(feed)

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
