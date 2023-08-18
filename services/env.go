package services

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AccessLifetime  int
	RefreshLifetime int
	Secret          string
	DbURI           string
	DbName          string
	DbTableName     string
	Port            string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessLifetime, err := strconv.Atoi(os.Getenv("ACCESS_LIFETIME"))
	if err != nil {
		fmt.Println("ACCESS_LIFETIME env variable can't be loaded")
	}

	refreshLifetime, err := strconv.Atoi(os.Getenv("REFRESH_LIFETIME"))
	if err != nil {
		fmt.Println("REFRESH_LIFETIME env variable can't be loaded")
	}

	fmt.Println()

	return Config{
		AccessLifetime:  accessLifetime,
		RefreshLifetime: refreshLifetime,
		Secret:          os.Getenv("SECRET_KEY"),
		DbURI:           os.Getenv("DB_URI"),
		DbName:          os.Getenv("DB_NAME"),
		DbTableName:     os.Getenv("DB_TABLE_NAME"),
		Port:            os.Getenv("PORT"),
	}
}
