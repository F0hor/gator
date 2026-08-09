package config

import(
	"os"
	"encoding/json"
	"log"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl string `json:"db_url"`
	User string `json:"current_user_name"`
}

func Read() Config {
	data, err := os.ReadFile(getConfigFilePath())
	if err != nil	{
		log.Fatal(err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func (c *Config) SetUser(newUser string) error {
	c.User = newUser
	return write(*c)
}

func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(getConfigFilePath(), data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return userHome + "/" + configFileName
}

