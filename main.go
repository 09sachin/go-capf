package main

import (
	"fmt"
	"github.com/09sachin/go-capf/routes"
	"github.com/09sachin/go-capf/controllers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}


func main() {

	initLoggers()
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

func initLoggers() {
	// Open log files
	infoLogFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open info log file:", err)
	}
	errorLogFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open error log file:", err)
	}

	controllers.InfoLogger = log.New(infoLogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	controllers.ErrorLogger = log.New(errorLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
