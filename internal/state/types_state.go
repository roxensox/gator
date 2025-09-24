package state

import (
	_ "github.com/lib/pq"
	"github.com/roxensox/gator/internal/config"
	"github.com/roxensox/gator/internal/database"
)

type State struct {
	Cfg  *config.Config
	Conn *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Reg map[string]func(s *State, cmd Command) error
}
