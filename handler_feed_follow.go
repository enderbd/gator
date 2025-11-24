package main

import (
	"context"
	"fmt"
	"time"

	"github.com/enderbd/gator/internal/database"
	"github.com/google/uuid"
)

func handleFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Missing argument, usage : %s <url> ", cmd.name)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Could not find a feed with the url %s: error %w", cmd.args[0], err)
	}

	currentUserName := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), currentUserName) 
	if err != nil {
		return fmt.Errorf("Could not find current user %v in the database: %+w ", currentUserName, err)
	}

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Could not create an entry into feed_follow")
	}

	fmt.Printf(" * Feed Name: %s\n", feed_follow.FeedName)
	fmt.Printf(" * User Name: %s\n", feed_follow.UserName)

	return nil
}


func handleFollowing(s *state, _ command) error {
	currentUserName := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("Error getting the current user %w", err)
	}

	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("No feed follows for user %s: error %w", currentUserName, err)
	}
	
	for _, feed_follow := range feed_follows {
		fmt.Printf("* %s\n", feed_follow.FeedName)
	}

	return nil
}
