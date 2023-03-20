package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	goEnv := os.Getenv("GO_ENV")
	fmt.Println(goEnv)
	if goEnv == "" || goEnv == "dev" {
		err := godotenv.Load(".env")
		if err != nil {
			return err
		}
	}

	return nil
}
