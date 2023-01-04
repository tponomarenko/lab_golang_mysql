package main

import (
	"fmt"
	"os"
)

type Settings struct {
	ServicePort string
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
		DbUsername:  os.Getenv("DB_USERNAME"),
		DbPassword:  os.Getenv("DB_PASSWORD"),
		DbPort:      os.Getenv("DB_PORT"),
		DbHost:      os.Getenv("DB_HOST"),
		DbName:      os.Getenv("DB_NAME"),
	}, nil
}
