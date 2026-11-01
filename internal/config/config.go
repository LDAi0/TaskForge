package config

import (
	"errors"
	"os"
)

type Config struct{
	DatabaseURL string
	HTTPaddress string
}

func Load() (Config, error){
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == ""{
		return Config{}, errors.New("Database URL is empty!")
	}

	httpaddress := os.Getenv("API_PORT")
	if httpaddress == ""{
		httpaddress = "8080"
	}

	return Config{
		DatabaseURL: databaseURL,
		HTTPaddress: httpaddress,
	}, nil
}