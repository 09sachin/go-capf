package controllers

import (
	"fmt"
	"encoding/json"
	"bytes"
	"time"
	"github.com/09sachin/go-capf/config"
	_"github.com/09sachin/go-capf/models"
	"net/http"
	"github.com/dgrijalva/jwt-go"
)


var jwtKey = []byte("your-secret-key")

// Claims structure to represent the data that will be encoded to create the JWT
type TokenClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}


type Response struct {
	Message string `json:"message"`
}

type JsonResponse struct {
	Message json.RawMessage `json:"message"`
}


type PhoneNo struct {
	MobileNumber    string
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err := validateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If the token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	}
}

func createToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &TokenClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateToken(tokenString string) (*TokenClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return claims, nil
}


func getClaimsFromRequest(r *http.Request) (*TokenClaim, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("No token provided")
	}

	claims, err := validateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("Invalid token: %v", err)
	}

	return claims, nil
}


func OtpLogin(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("force_id")
	login_q := fmt.Sprintf(`select mobile_number, relation_name, id_number  
	from capf.capf_prod_noimage_refresh 
	where id_number='%s' and relation_name='Self'`, id)
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
	id := r.FormValue("force_id")
	login_q := fmt.Sprintf(`select mobile_number  
	from capf.capf_prod_noimage_refresh 
	where id_number='%s' and relation_name='Self'`, id)

	rows, _ := config.ExecuteQuery(login_q)
	
	var dataList []PhoneNo

	for rows.Next() {
		var data PhoneNo
		err := rows.Scan(&data.MobileNumber)
		fmt.Println(err)
		dataList = append(dataList, data)	
	}

	phone_no := dataList[0].MobileNumber

	success := OtpService((phone_no))

	if success==0{
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func OtpService(phone string) int{
	// URL to which the POST request will be sent
	url := "https://example.com/api"

	// Request payload (data to be sent in the request body)
	payload := []byte(fmt.Sprintf(`{"phone_no": "%s", "key2": "value2"}`, phone))

	// Create a new POST request with the specified URL and payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 0
	}

	// Set the Content-Type header for the request
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return 0
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode == http.StatusOK {
		fmt.Println("POST request was successful")
		return 1
	} else {
		fmt.Println("POST request failed with status:", resp.Status)
		return 0
	}
}
