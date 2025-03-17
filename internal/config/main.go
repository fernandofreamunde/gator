package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fernandofreamunde/gator/internal/database"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	Db     *database.Queries
	Config *Config
}

func Read() Config {
	fmt.Println("reading...")
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.ReadFile(home + "/.gatorconfig.json")
	if err != nil {
		fmt.Println(err)
	}

	config := Config{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(config)
	return config
}

func (c *Config) SetUser(username string) error {

	c.CurrentUserName = username

	fmt.Println("writing...")
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
	}

	err = os.WriteFile(home+"/.gatorconfig.json", bytes, 0644)
	if err != nil {
		fmt.Println(err)
	}

	//	file, err := os.ReadFile(home + "/.gatorconfig.json")
	//	if err != nil {
	//		fmt.Println(err)
	//	}

	return nil
}
