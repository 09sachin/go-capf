package main

import (
	"fmt"
	"errors"
	"github.com/ichtrojan/go-todo/routes"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	
	if err := godotenv.Load(); err != nil {
		fmt.Println(errors.New("no .env file found"))
	}

	port, exist := os.LookupEnv("PORT")

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
