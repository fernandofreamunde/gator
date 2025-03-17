package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/fernandofreamunde/gator/internal/commands"
	"github.com/fernandofreamunde/gator/internal/config"
	"github.com/fernandofreamunde/gator/internal/database"
)

func main() {
	fmt.Println("Hello Gator!")
	s := config.State{}
	c := config.Read() // maybe return errors? and detect and exit with os.Exit(1)
	s.Config = &c

	db, err := sql.Open("postgres", c.DbUrl)
	if err != nil {
		fmt.Println("Error: Failed to connect to the database:", err)
	}

	dbQueries := database.New(db)
	s.Db = dbQueries

	cmds := &commands.Commands{
		Registry: make(map[string]func(*config.State, commands.Command) error),
	}

	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerUsers)

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

	err = cmds.Run(&s, cmd)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func handlerLogin(s *config.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is requiered")
	}

	ctx := context.Background()
	name := sql.NullString{String: cmd.Args[0], Valid: true}
	existingUser, _ := s.Db.GetUserByName(ctx, name)
	if !existingUser.Name.Valid {
		return fmt.Errorf("user does not exist")
	}
	err := s.Config.SetUser(cmd.Args[0])

	if err != nil {
		return err
	}
	fmt.Println("user has been sucessfuly set")

	return nil
}

func handlerRegister(s *config.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is requiered")
	}
	ctx := context.Background()
	name := sql.NullString{String: cmd.Args[0], Valid: true}
	existingUser, _ := s.Db.GetUserByName(ctx, name)
	if existingUser.Name.Valid {
		return fmt.Errorf("user already exists")
	}

	user, err := s.Db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	fmt.Println("user has been sucessfuly created")

	err = s.Config.SetUser(user.Name.String)
	if err != nil {
		return err
	}
	fmt.Println("user has been sucessfuly set")

	return nil
}

func handlerReset(s *config.State, cmd commands.Command) error {
	ctx := context.Background()

	err := s.Db.ResetUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Println("users have been sucessfuly reset")

	return nil
}

func handlerUsers(s *config.State, cmd commands.Command) error {
	ctx := context.Background()

	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if s.Config.CurrentUserName == user.Name.String {

			fmt.Println("* ", user.Name.String, "(current)")
			continue
		}
		fmt.Println("* ", user.Name.String)
	}
	return nil
}
