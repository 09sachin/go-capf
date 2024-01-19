package controllers

import (
	"fmt"
	"encoding/json"
	_"github.com/09sachin/go-capf/config"
	_"github.com/09sachin/go-capf/models"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

type JsonResponse struct {
	Message json.RawMessage `json:"message"`
}

func OtpLogin(w http.ResponseWriter, r *http.Request) {
	login_q := `select mobile_number, relation_name, id_number  
	from capf.capf_prod_noimage_refresh 
	where id_number='000000523' and relation_name='Self'`
	fmt.Println(login_q)
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
