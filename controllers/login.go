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
	PmjayId string `json:"pmjayid"`
	ForceType string `json:"force_type"`
	jwt.StandardClaims
}


type Response struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type JsonResponse struct {
	Message json.RawMessage `json:"message"`
}


type PhoneNo struct {
	MobileNumber    string
}

type Pmjay struct {
	PmjayId    string
}

type OTP struct {
	Otp    	   string
	Created_at time.Time   
}

// func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tokenString := r.Header.Get("Authorization")

// 		if tokenString == "" {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		_, err := validateToken(tokenString)
// 		if err != nil {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// If the token is valid, proceed to the next handler
// 		next.ServeHTTP(w, r)
// 	}
// }

func createToken(username string, pmjay string, force_type string) (string, error) {
	expirationTime := time.Now().Add(120 * time.Minute)

	claims := &TokenClaim{
		Username: username,
		PmjayId: pmjay,
		ForceType: force_type,
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
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}


func getClaimsFromRequest(r *http.Request) (*TokenClaim, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("no token provided")
	}

	claims, err := validateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return claims, nil
}

type RequestBody struct {
	ForceID 	string `json:"force_id"`
	OTP     	string `json:"otp"`
	ForceType   string `json:"force_type"`
}

type PmjayQuery struct{
	PMJAY string
}


func formatStringSlice(slice []string) string {
	result := ""
	for i, value := range slice {
		result += fmt.Sprintf("'%s'", value)
		if i < len(slice)-1 {
			result += ", "
		}
	}
	return result
}


func OtpLogin(w http.ResponseWriter, r *http.Request) {
	/// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Error reading request body",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Unmarshal the JSON data into a struct
	var requestData RequestBody
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Error decoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Access the values from the struct
	id := requestData.ForceID
	otp := requestData.OTP
	force_type := requestData.ForceType
	login_id := force_type + "-" + id
	fmt.Println(login_id)
	get_otp := fmt.Sprintf("select otp, updated_at from login where force_id='%s'", login_id)

	rows, err := config.ExecuteQueryLocal(get_otp)
	fmt.Println((err))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var dataList []OTP
	for rows.Next() {
		fmt.Println("yes")
		var data OTP
		err := rows.Scan(&data.Otp, &data.Created_at)
		if err!=nil{
			fmt.Println(err)
		}
		dataList = append(dataList, data)	
	}

	otp_stored := dataList[0].Otp
	exp_time := dataList[0].Created_at.Add(10 * time.Minute) 

	// otp_stored := "123456"

	if otp_stored!=otp{
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Incorrect OTP",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	current_time := (time.Now().UTC().Add(330*time.Minute))
	if exp_time.Before(current_time){
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "OTP expired, please resend OTP",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	pmjay_q := fmt.Sprintf(`select distinct pmjay_id  
	from capf.capf_prod_noimage_refresh 
	where id_number='%s' and id_type='%s'`, id, force_type)

	rows, sql_error := config.ExecuteQuery(pmjay_q)
	if sql_error!=nil{
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var pmjayids []string
	for rows.Next() {
		var pmjayid string
		err := rows.Scan(&pmjayid)
		if err != nil {
			return
		}
		pmjayids = append(pmjayids, pmjayid)
	}

	str := "(" + formatStringSlice(pmjayids) + ")"
	
	token, _ := createToken(id, str, force_type)

	response  := Response{
		Message:  token,
	}

	
	// Encode the response as JSON and write it to the response writer
	err2 := json.NewEncoder(w).Encode(response)
	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Error encoding JSON",
		}
		json.NewEncoder(w).Encode(response)
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
	w.Header().Set("Content-Type", "application/json")

	// Read the request body
	body, err1 := ioutil.ReadAll(r.Body)
	if err1 != nil {
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Error reading request body",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Unmarshal the JSON data into a struct
	var requestData RequestBody
	err1 = json.Unmarshal(body, &requestData)
	if err1 != nil {
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Error decoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Access the values from the struct
	id := requestData.ForceID
	force_type := requestData.ForceType
	login_id := force_type + "-" + id

	login_q := fmt.Sprintf(`select mobile_number  
	from capf.capf_prod_noimage_refresh 
	where id_number='%s' and id_type='%s' and relation_name='Self'`, id, force_type)

	rows, sql_error := config.ExecuteQuery(login_q)
	if sql_error!=nil{
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var dataList []PhoneNo

	for rows.Next() {
		var data PhoneNo
		err := rows.Scan(&data.MobileNumber)
		if err!=nil{
			fmt.Println(err)
		}
		dataList = append(dataList, data)	
	}

	if len(dataList)==0{
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Wrong force id",
		}
		json.NewEncoder(w).Encode(response)
		return
	}


	// phone_og := dataList[0].MobileNumber
	phone_og := "7014600922"
	otp := generateOTP()
	//otp := "123456"
	save_otp_query := fmt.Sprintf(`INSERT INTO login (force_id, otp)
	VALUES ('%s', '%s')
	ON CONFLICT (force_id)
	DO UPDATE SET otp = %s;`, login_id, otp, otp)
	row:= config.InsertData(save_otp_query)
	fmt.Println(row)
	

	success := sendSMSAPI(phone_og, otp)
	var message string
	if success{
		message = fmt.Sprintf("OTP sent successfully to %s", phone_og)
	}else{
		message = fmt.Sprintf("Failed to send OTP to %s", phone_og)
	}

	if !success{
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Failed to send OTP, please try again",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response  := Response{
		Message:  message,
	}

	// Encode the response as JSON and write it to the response writer
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response  := ErrorResponse{
			Error:  "Error encoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func sendSMSAPI(phoneNo, otp string) bool {
	msg := "Dear%20User%2C%0AYour%20OTP%20to%20access%20CAPF%20application%20is%20ABCDEF.%20It%20will%20be%20valid%20for%203%20minutes.%0ANHA"
	msg = strings.Replace(msg, "ABCDEF", otp, -1)
	// fmt.Println(msg)
	username := "abhaotp"
	password := "f9F3r%5D%7BS"
	entityID := "1001548700000010184"
	tempID := "1007170748130898041"
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
