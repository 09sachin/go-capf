package main

import (
	"fmt"
	"errors"
	"github.com/09sachin/go-capf/routes"
	"github.com/joho/godotenv"
	"net/http"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {

	loadEnv()

	port, exist := "8000", true

	if !exist {
		fmt.Println(errors.New("PORT not set in .env"))
		fmt.Println("PORT not set in .env")
	}

	err := http.ListenAndServe(":"+port, routes.Init())

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
}
