package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/enderbd/gator/internal/database"
)


func handleBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	if len(cmd.args) < 1 {
		limit = 2
	} else {
		convertedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("When provided, the argument must be an integer. Usage: %s <limit>", cmd.name)
		}
		limit = int32(convertedLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: limit,
	})
	if err != nil {
		return fmt.Errorf("Error retrieving the post: %w", err)
	}

	for _, post := range posts {
		fmt.Printf(" Post title: %s\n", post.Title.String)
		fmt.Printf(" Post: %s\n", post.Description.String)
	}
	return nil
	

}
