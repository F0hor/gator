package main

import(
	"fmt"
	"context"

	"github.com/F0hor/database"
)

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.User)
		if err != nil {
			return fmt.Errorf("Failed to retieve nessery data:\n%v\n", err)
		}

		return handler(s, cmd, user)
	}
}

