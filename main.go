package main

import (
	"fmt"
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
		fmt.Println("PORT not set in .env")
	}

	err := http.ListenAndServe(":"+port, routes.Init())

	if err != nil {
		fmt.Println(err)
	}
}
