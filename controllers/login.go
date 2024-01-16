package controllers

import (
	"fmt"
	"encoding/json"
	// "github.com/09sachin/go-capf/config"
	// "github.com/09sachin/go-capf/models"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func OtpLogin(w http.ResponseWriter, r *http.Request) {

	response := Response{
		Message: "Hello, JSON!",
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response writer
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}


func SendOtp(w http.ResponseWriter, r *http.Request) {

	fmt.Println("opt sent")

	http.Redirect(w, r, "/", 301)
}
