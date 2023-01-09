package main

import (
	"fmt"
	"os"
)

type Settings struct {
	ServicePort string
	AuthToken   string
	DbEngine    string
	DbUsername  string
	DbPassword  string
	DbPort      string
	DbHost      string
	DbName      string
}

func NewSettings() (*Settings, error) {
	servicePort, exists := os.LookupEnv("SERVICE_PORT")

	if !exists {
		msg := "SERVICE_PORT is not specified"
		fmt.Println(msg)
		return nil, fmt.Errorf(msg)
	}

	return &Settings{
		ServicePort: servicePort,
		AuthToken:   os.Getenv("AUTH_TOKEN"),
		DbUsername:  os.Getenv("DB_USERNAME"),
		DbEngine:    os.Getenv("DB_ENGINE"),
		DbPassword:  os.Getenv("DB_PASSWORD"),
		DbPort:      os.Getenv("DB_PORT"),
		DbHost:      os.Getenv("DB_HOST"),
		DbName:      os.Getenv("DB_NAME"),
	}, nil
}
