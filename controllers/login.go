package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/09sachin/go-capf/config"
	"github.com/09sachin/go-capf/models"
	"html/template"
	"net/http"
)

var (
	id        int
	item      string
	completed int
	database  = config.Database()
)

func Login(w http.ResponseWriter, r *http.Request) {

	fmt.Println("login")

	http.Redirect(w, r, "/", 301)
}


func SendOtp(w http.ResponseWriter, r *http.Request) {

	fmt.Println("opt sent")

	http.Redirect(w, r, "/", 301)
}
