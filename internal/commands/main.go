package commands

import (
	"fmt"

	"github.com/fernandofreamunde/gator/internal/config"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Registry map[string]func(*config.State, Command) error
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
	c.Registry[name] = f
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	command, ok := c.Registry[cmd.Name]
	if !ok {
		return fmt.Errorf("Unknown command: %s", cmd.Name)
	}

	return command(s, cmd)
}
