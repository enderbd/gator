package main

import (
	"context"
	"fmt"
	"time"

	"github.com/enderbd/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command ) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("username missing, usage: %s <username>", cmd.name)
	}
	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("could not find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("Coiuld not set the user name: %w", err)
	}

	fmt.Println("Username set!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("username missing, usage: %s <username>", cmd.name)
	}

	name := cmd.args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams {
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

func handleReset(s *state, _ command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Could not reset the database!")
	}
	fmt.Println("Database reseted, success!")
	return nil
}

func handleUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error trying to get all the users: %w", err)	
	}

	for _, user := range users {
		name := user.Name
		if name == s.cfg.CurrentUserName {
			name += " (current)"
		}
		fmt.Printf("* %s\n", name)
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
