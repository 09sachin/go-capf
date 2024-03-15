package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	_ "net/url"
	"strconv"
	"strings"
	"time"
	"github.com/09sachin/go-capf/config"
	_ "github.com/09sachin/go-capf/models"
)


func OtpLogin(w http.ResponseWriter, r *http.Request) {
	/// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error reading request body",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Unmarshal the JSON data into a struct
	var requestData RequestBody
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error decoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Access the values from the struct
	id := requestData.ForceID
	otp := requestData.OTP
	force_type := requestData.ForceType
	login_id := force_type + "-" + id
	InfoLogger.Println(login_id)
	get_otp := fmt.Sprintf("select otp, updated_at from login where force_id='%s'", login_id)

	rows, err := config.ExecuteQueryLocal(get_otp)
	fmt.Println(err)
	if err != nil {
		ErrorLogger.Printf("Database connection error : %s", login_id)
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var dataList []OTP
	for rows.Next() {
		var data OTP
		err := rows.Scan(&data.Otp, &data.Created_at)
		if err != nil {
			fmt.Println(err)
		}
		dataList = append(dataList, data)
	}

	if len(dataList) == 0 {
		ErrorLogger.Printf("Unauthorised access to otp login : %s", login_id)
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Please send otp again",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	otp_stored := dataList[0].Otp
	exp_time := dataList[0].Created_at.Add(10 * time.Minute)

	// otp_stored := "123456"

	if otp_stored != otp {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Incorrect OTP",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	current_time := (time.Now().UTC().Add(330 * time.Minute))
	if exp_time.Before(current_time) {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "OTP expired, please resend OTP",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	pmjay_q := fmt.Sprintf(`select distinct pmjay_id  
	from capf_prod_noimage_refresh 
	where id_number='%s' and id_type='%s'`, id, force_type)

	rows, sql_error := config.ExecuteQuery(pmjay_q)
	if sql_error != nil {
		ErrorLogger.Printf("Database connection error : %s", login_id)
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Database connection could not be established",
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

	response := Response{
		Message: token,
	}

	// Encode the response as JSON and write it to the response writer
	err2 := json.NewEncoder(w).Encode(response)
	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error encoding JSON",
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
	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error reading request body",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Unmarshal the JSON data into a struct
	var requestData RequestBody
	err1 = json.Unmarshal(body, &requestData)
	if err1 != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error decoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Access the values from the struct
	id := requestData.ForceID
	force_type := requestData.ForceType
	login_id := force_type + "-" + id
	InfoLogger.Println(login_id)

	login_q := fmt.Sprintf(`select mobile_number  
	from capf_prod_noimage_refresh 
	where id_number='%s' and id_type='%s' and relation_name='Self'`, id, force_type)

	rows, sql_error := config.ExecuteQuery(login_q)
	if sql_error != nil {
		ErrorLogger.Printf("Database connection error : %s", login_id)
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var dataList []PhoneNo

	for rows.Next() {
		var data PhoneNo
		err := rows.Scan(&data.MobileNumber)
		if err != nil {
			fmt.Println(err)
		}
		dataList = append(dataList, data)
	}

	if len(dataList) == 0 {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Wrong force id",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// phone_og := dataList[0].MobileNumber
	phone_og := "6377035564"
	otp := generateOTP()
	// otp := "123456"
	save_otp_query := fmt.Sprintf(`INSERT INTO login (force_id, otp)
	VALUES ('%s', '%s')
	ON CONFLICT (force_id)
	DO UPDATE SET otp = %s;`, login_id, otp, otp)
	config.InsertData(save_otp_query)

	success := sendSMSAPI(phone_og, otp)
	var message string
	if success {
		message = fmt.Sprintf("OTP sent successfully to %s", phone_og)
	} else {
		message = fmt.Sprintf("Failed to send OTP to %s", phone_og)
	}

	// if !success{
	// 	w.WriteHeader(http.StatusNotFound)
	// 	response  := ErrorResponse{
	// 		Error:  "Failed to send OTP, please try again",
	// 	}
	// 	json.NewEncoder(w).Encode(response)
	// 	return
	// }

	response := Response{
		Message: message,
	}

	// Encode the response as JSON and write it to the response writer
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error encoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func sendSMSAPI(phoneNo, otp string) bool {
	msg := "Dear%20User%2C%0AYour%20OTP%20to%20access%20CAPF%20application%20is%20ABCDEF.%20It%20will%20be%20valid%20for%203%20minutes.%0ANHA"
	msg = strings.Replace(msg, "ABCDEF", otp, -1)
	username := "abhaotp"
	password := "f9F3r%5D%7BS"
	entityID := "1001548700000010184"
	tempID := "1007170748130898041"
	source := "NHASMS"

	urlStr := fmt.Sprintf("https://sms6.rmlconnect.net/bulksms/bulksms?username=%s&password=%s&type=0&dlr=1&destination=%s&source=%s&message=%s&entityid=%s&tempid=%s",
		username, password, phoneNo, source, msg, entityID, tempID)

	response, err := http.Post(urlStr, "application/json", nil)
	if err != nil {
		ErrorLogger.Printf("SMS API failed")
		return false
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("SMS sent successfully")
		return true
	}

	ErrorLogger.Println("Failed to send SMS. Status code : ", response.StatusCode)
	return false
}
