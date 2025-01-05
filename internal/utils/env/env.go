package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv("APP_ENV")

	if env != "" {
		godotenv.Load(fmt.Sprintf(".env.%s.local", env))
		if env != "test" {
			godotenv.Load(fmt.Sprintf(".env.%s", env))
		}
	}

	godotenv.Load(".env.local")
	godotenv.Load(".env")
}
