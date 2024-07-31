package static

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var Config = struct {
	KeyPair struct {
		AccessKey string
		SecretKey string
	}
}{}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}
	Config.KeyPair.AccessKey = os.Getenv("ACCESS_KEY")
	Config.KeyPair.SecretKey = os.Getenv("SECRET_KEY")
}
