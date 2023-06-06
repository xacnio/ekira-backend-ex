package configs

import (
	"github.com/joho/godotenv"
	"os"
)

// LoadConfig func for environment variables of .env files
func LoadConfig() error {
	if os.Getenv("RUN_TYPE") == "" || os.Getenv("RUN_TYPE") == "dev" {
		err := godotenv.Load(".env.dev")
		if err != nil {
			return err
		}
	} else if os.Getenv("RUN_TYPE") == "PROD" {
		err := godotenv.Load(".env.prod")
		if err != nil {
			return err
		}
	}
	return nil
}
