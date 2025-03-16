package main

import (
	"fmt"
	"os"

	"github.com/fernandofreamunde/gator/internal/commands"
	"github.com/fernandofreamunde/gator/internal/config"
)

func main() {
	fmt.Println("Hello Gator!")
	s := config.State{}
	c := config.Read() // maybe return errors? and detect and exit with os.Exit(1)
	s.Config = &c

	cmds := &commands.Commands{
		Registry: make(map[string]func(*config.State, commands.Command) error),
	}

	cmds.Register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Printf("Error: not enough arguments")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := commands.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err := cmds.Run(&s, cmd)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func handlerLogin(s *config.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is requiered")
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("user has been sucessfuly set")

	return nil
}
