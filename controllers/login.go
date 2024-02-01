package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	_"net/url"
	"strconv"
	"strings"
	"time"
	"github.com/09sachin/go-capf/config"
	_ "github.com/09sachin/go-capf/models"
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

type OTP struct {
	Otp    string
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
	expirationTime := time.Now().Add(120 * time.Minute)

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

type RequestBody struct {
	ForceID string `json:"force_id"`
	OTP     string `json:"otp"`
}

type PmjayQuery struct{
	PMJAY string
}

func OtpLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a struct
	var requestData RequestBody
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Access the values from the struct
	id := requestData.ForceID
	otp := requestData.OTP

	get_otp := fmt.Sprintf("select otp from login where force_id='%s'", id)

	rows, err := config.ExecuteQueryLocal(get_otp)
	if err!=nil{
		fmt.Println(err)
		return
	}
	var dataList []OTP

	for rows.Next() {
		var data OTP
		err := rows.Scan(&data.Otp)
		fmt.Println(err)
		dataList = append(dataList, data)	
	}

	otp_stored := dataList[0].Otp

	// otp_stored := "123456"
	fmt.Println(otp)

	if otp_stored!=otp{
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, _ := createToken(id)

	response  := Response{
		Message:  token,
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response writer
	err2 := json.NewEncoder(w).Encode(response)
	if err2 != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}


func generateOTP() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	otp := r.Intn(900000) + 100000
	otp_str := strconv.Itoa(otp)
	return otp_str
}

func SendOtp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err1 := ioutil.ReadAll(r.Body)
	if err1 != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a struct
	var requestData RequestBody
	err1 = json.Unmarshal(body, &requestData)
	if err1 != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Access the values from the struct
	id := requestData.ForceID
	fmt.Println(id)
	login_q := fmt.Sprintf(`select mobile_number  
	from capf.capf_prod_noimage_refresh 
	where id_number='%s' and relation_name='Self'`, id)

	rows, sql_error := config.ExecuteQuery(login_q)
	if sql_error!=nil{
		fmt.Println(sql_error)
		return
	}
	var dataList []PhoneNo

	for rows.Next() {
		var data PhoneNo
		err := rows.Scan(&data.MobileNumber)
		fmt.Println(err)
		dataList = append(dataList, data)	
	}
	fmt.Println(len(dataList))

	if len(dataList)==0{
		http.Error(w, "Wrong force id", http.StatusNotFound)
		return
	}


	phone_og := dataList[0].MobileNumber
	phone_no := "6377035564"
	// otp := generateOTP()
	otp := "123456"
	save_otp_query := fmt.Sprintf(`INSERT INTO login (force_id, otp)
	VALUES ('%s', '%s')
	ON CONFLICT (force_id)
	DO UPDATE SET otp = %s;`, id, otp, otp)
	row:= config.InsertData(save_otp_query)
	fmt.Println(row)
	
	fmt.Println(otp)

	success := sendSMSAPI(phone_no, otp)
	var message string
	if success{
		message = fmt.Sprintf("OTP sent successfully to %s", phone_og)
	}else{
		message = fmt.Sprintf("Failed to send OTP to %s", phone_og)
	}

	response  := Response{
		Message:  message,
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

func sendSMSAPI(phoneNo, otp string) bool {
	msg := "Dear%20User%2C%0AOTP%20to%20validate%20your%20Allied%20and%20Healthcare%20Institute%20Registry%20application%20is%20ABCDEF.%20This%20is%20One%20Time%20Password%20will%20be%20valid%20for%2010%20mins.%0AABDM%2C%20National%20Health%20Authority"
	msg = strings.Replace(msg, "ABCDEF", otp, -1)
	// fmt.Println(msg)
	username := "abhaotp"
	password := "f9F3r%5D%7BS"
	entityID := "1001548700000010184"
	tempID := "1007169865765792689"
	source := "NHASMS"

	urlStr := fmt.Sprintf("https://sms6.rmlconnect.net/bulksms/bulksms?username=%s&password=%s&type=0&dlr=1&destination=%s&source=%s&message=%s&entityid=%s&tempid=%s",
		username, password, phoneNo, source, msg, entityID, tempID)
	
	response, err := http.Post(urlStr, "application/json", nil)
	if err != nil {
		fmt.Println("Error sending SMS:", err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("SMS sent successfully")
		return true
	}

	fmt.Println("Failed to send SMS. Status code:", response.StatusCode)
	return false
}
