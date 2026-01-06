package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db    DbConfig
	Auth  AuthConfig
	Token TokenConfig
}
type DbConfig struct {
	Dsn string
}
type AuthConfig struct {
	Auth string
}
type TokenConfig struct {
	AdminToken string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error with loading config")
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			os.Getenv("AUTH"),
		},
		Token: TokenConfig{
			os.Getenv("ADMINTOKEN"),
		},
	}
}
