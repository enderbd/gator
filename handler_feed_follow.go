package main

import (
	"context"
	"fmt"
	"time"

	"github.com/enderbd/gator/internal/database"
	"github.com/google/uuid"
)

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Missing argument, usage : %s <feed_url> ", cmd.name)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Could not find a feed with the url %s: error %w", cmd.args[0], err)
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


func handleFollowing(s *state, _ command, user database.User) error {
	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("No feed follows for user %s: error %w", user.Name, err)
	}
	
	for _, feed_follow := range feed_follows {
		fmt.Printf("* %s\n", feed_follow.FeedName)
	}

	return nil
}

func handleUnfollow(s  *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Missing argument, usage : %s <feed_url> ", cmd.name)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Could not find feed with url %s to unfollow: error %w", cmd.args[0], err)
	}

	if err := s.db.DeleteFeedFollowByUserFeedIds(context.Background(), database.DeleteFeedFollowByUserFeedIdsParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("Unable to delete feed with url %s from feed_follows: error %w", cmd.args[0], err)
	}

	fmt.Printf("Followed feed removed!")
	return nil
}
