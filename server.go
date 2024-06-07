package main

import (
	"os"

	"library/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		return
	}

	AppPort := os.Getenv("APP_PORT")

	e := routes.Init()

	e.Logger.Fatal(e.Start(":" + AppPort))
}
