package state

import (
	"context"
	"github.com/roxensox/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Conn.GetUser(context.Background(), *s.Cfg.CurrentUser)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
