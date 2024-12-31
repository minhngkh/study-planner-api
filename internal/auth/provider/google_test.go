package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func LoadENV() {
	path := ".env"
	for {
		err := godotenv.Load(path)
		if err == nil {
			break
		}
		path = "../" + path
	}
}

func TestGoogle(t *testing.T) {
	LoadENV()

	// GoogleLogin()
	fmt.Println(os.Getenv("GOOGLE_CLIENT_ID"))
	fmt.Println("TestGoogle")
}
