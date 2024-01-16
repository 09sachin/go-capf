package controllers

import (
	"fmt"
	// "github.com/09sachin/go-capf/config"
	// "github.com/09sachin/go-capf/models"
	"net/http"
)


func DashboardData(w http.ResponseWriter, r *http.Request) {

	fmt.Println("login")

	http.Redirect(w, r, "/", 301)
}


func UserDetails(w http.ResponseWriter, r *http.Request) {

	fmt.Println("opt sent")

	http.Redirect(w, r, "/", 301)
}



func Hospitals(w http.ResponseWriter, r *http.Request) {

	fmt.Println("opt sent")

	http.Redirect(w, r, "/", 301)
}

func FilterHospital(w http.ResponseWriter, r *http.Request) {

	fmt.Println("opt sent")

	http.Redirect(w, r, "/", 301)
}


func Queries(w http.ResponseWriter, r *http.Request) {

	fmt.Println("opt sent")

	http.Redirect(w, r, "/", 301)
}


func TrackCases(w http.ResponseWriter, r *http.Request) {

	fmt.Println("opt sent")

	http.Redirect(w, r, "/", 301)
}

func Claims(w http.ResponseWriter, r *http.Request) {

	fmt.Println("opt sent")

	http.Redirect(w, r, "/", 301)
}