package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DBString  = ""
	Port      = 0
	SecretKey []byte
)

func LoadConfig() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 9000
	}

	DBString = os.Getenv("DB_URL")

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
