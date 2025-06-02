package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/09sachin/go-capf/config"
	_ "github.com/09sachin/go-capf/models"
	"io"
	"math/rand"
	"net/http"
	_ "net/url"
	"strconv"
	"strings"
	"time"
)

func OtpLogin(w http.ResponseWriter, r *http.Request) {
	/// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		JsonParseError(w)
		return
	}

	// Unmarshal the JSON data into a struct
	var requestData RequestBody
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		JsonParseError(w)
		return
	}

	// Access the values from the struct
	id := requestData.ForceID
	otp := requestData.OTP
	force_type := requestData.ForceType
	isvalid := (isAlphaNumeric(id) && isAlphaNumeric(otp) && isAlphaNumeric(force_type))
	if !isvalid{
		QueryParamsError(w)
		return
	}

	phone_og := "916377035564"
	success := sendSMSAPInew(phone_og, otp)
	var message string

	masked_phone := maskPhoneNumber(phone_og)
	
	if success {
		message = fmt.Sprintf("OTP sent successfully to %s", masked_phone)
	} else {
		message = fmt.Sprintf("Failed to send OTP to %s", masked_phone)
	}
	Custom4O4Error(w,message)
	return
	InfoLogger.Println(message)

	
	login_id := force_type + "-" + id
	InfoLogger.Println(login_id)
	get_otp := "select otp, updated_at from login where force_id=$1"

	rows, err := config.ExecuteQueryLocal(get_otp, login_id)

	if err != nil {
		DbError(w)
		return
	}
	var dataList []OTP
	for rows.Next() {
		var data OTP
		err := rows.Scan(&data.Otp, &data.Created_at)
		if err != nil {
			ErrorLogger.Println(err)
		}
		dataList = append(dataList, data)
	}

	if len(dataList) == 0 {
		ErrorLogger.Printf("Unauthorised access to otp login : %s", login_id)
		// save_otp_query := `INSERT INTO login (force_id, otp)
		// VALUES ($1, $2)
		// ON CONFLICT (force_id)
		// DO UPDATE SET otp = $2;`
		// config.InsertData(save_otp_query, login_id, "701460")
		Custom4O4Error(w,"OTP Error")
		return
	}

	otp_stored := dataList[0].Otp
	exp_time := dataList[0].Created_at.Add(10 * time.Minute)

	// default_otp := "701460"

	if (otp_stored != otp) {
		Custom4O4Error(w,"Incorrect OTP")
		return
	}

	current_time := (time.Now().UTC().Add(330 * time.Minute))
	if (exp_time.Before(current_time)){
		Custom4O4Error(w,"OTP expired, please resend OTP")
		return
	}

	//Test conditon for playstore/appstore
	if (id=="00000000" && force_type=="BS"){
		xtoken, _ := createToken(id, "('test')", "name_string", force_type)
		responsex := Response{
			Message: xtoken,
		}
	
		errx := json.NewEncoder(w).Encode(responsex)
		if errx != nil {
			JsonEncodeError(w)
		}
		return
	}
	

	urlStr := Beneficiary_URL
	payload := map[string]string{
		"id_type":   force_type,
		"id_number": id,
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		JsonEncodeError(w)
		return
	}
	search_response, err := http.Post(urlStr, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		ErrorLogger.Printf("Search API failed")
		ErrorLogger.Println(err)
		Custom4O4Error(w,"Search request failed")
		return
	}

	defer search_response.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(search_response.Body).Decode(&result)
	if err != nil {
		JsonParseError(w)
		return
	}
	if result["details"]==nil{
		Custom4O4Error(w,"Wrong force id / request failed")
		return
	}
	detailsArray := result["details"].([]interface{})

	var pmjayids []string
	var names []string

	for _, item := range detailsArray {
		detail := item.(map[string]interface{})
		pmjayids = append(pmjayids, detail["pmjay_id"].(string)) 
		names = append(names, detail["member_name_eng"].(string))
	}

	str := "(" + formatStringSlice(pmjayids) + ")"
	names_string := formatStringSlice(names)

	token, _ := createToken(id, str, names_string, force_type)

	response := Response{
		Message: token,
	}

	// Encode the response as JSON and write it to the response writer
	err2 := json.NewEncoder(w).Encode(response)
	if err2 != nil {
		JsonEncodeError(w)
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
		JsonParseError(w)
		return
	}

	// Unmarshal the JSON data into a struct
	var requestData RequestBody
	err1 = json.Unmarshal(body, &requestData)
	if err1 != nil {
		JsonParseError(w)
		return
	}

	// Access the values from the struct
	id := requestData.ForceID
	force_type := requestData.ForceType
	login_id := force_type + "-" + id
	InfoLogger.Println(login_id)

	//Test conditon for playstore/appstore
	if (login_id=="BS-00000000"){
		save_otp_query := `INSERT INTO login (force_id, otp)
		VALUES ($1, $2)
		ON CONFLICT (force_id)
		DO UPDATE SET otp = $2;`
		config.InsertData(save_otp_query, login_id, "123456")
		message_x := "OTP sent successfully to XXXXXX0000"
		response_x := Response{
			Message: message_x,
		}
		_ = json.NewEncoder(w).Encode(response_x)
		return
	}

	urlStr := Beneficiary_URL
	// Create JSON payload
	payload := map[string]string{
		"id_type":   force_type,
		"id_number": id,
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		JsonParseError(w)
		return
	}
	search_response, err := http.Post(urlStr, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		ErrorLogger.Printf("Search API failed")
		ErrorLogger.Println(err)
		Custom4O4Error(w,"Search request failed")
		return
	}

	defer search_response.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(search_response.Body).Decode(&result)
	if err != nil {
		JsonParseError(w)
		return
	}
	if result["details"]==nil{
		Custom4O4Error(w,"Wrong force id / request failed")
		return
	}
	detailsArray := result["details"].([]interface{})

	var self_data map[string]interface{}

	for _, item := range detailsArray {
		// Convert the item to a map[string]interface{}
		detail := item.(map[string]interface{})

		// Check if the "member_type" is "S"
		if detail["member_type"] == "S" {
			self_data = detail
		}
	}

	if search_response.StatusCode != http.StatusOK {
		Custom4O4Error(w,"Wrong force id / request failed")
		return
	}

	var phone_og string
	if mobileNumber, ok := self_data["mobile_number"].(string); ok {
		phone_og = mobileNumber
		isvalid := isAlphaNumeric(phone_og)
		if !isvalid{
			QueryParamsError(w)
			return
		}
	} else {
		Custom4O4Error(w, "Phone no not found")
		return
	}
	
	otp := generateOTP()
	// otp := "123456"
	save_otp_query := `INSERT INTO login (force_id, otp)
	VALUES ($1, $2)
	ON CONFLICT (force_id)
	DO UPDATE SET otp =$2;`
	config.InsertData(save_otp_query, login_id, otp)

	success := sendSMSAPInew(phone_og, otp)
	var message string

	masked_phone := maskPhoneNumber(phone_og)
	
	if success {
		message = fmt.Sprintf("OTP sent successfully to %s", masked_phone)
	} else {
		message = fmt.Sprintf("Failed to send OTP to %s", masked_phone)
	}

	if !success{
		Custom4O4Error(w,"Failed to send OTP, please try again")
		return
	}

	response := Response{
		Message: message,
	}

	// Encode the response as JSON and write it to the response writer
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		JsonEncodeError(w)
		return
	}
}

func sendSMSAPI(phoneNo, otp string) bool {
	if otp=="123456"{
		return true
	}
	msg := "Dear%20User%2C%0AYour%20OTP%20to%20access%20CAPF%20application%20is%20ABCDEF.%20It%20will%20be%20valid%20for%203%20minutes.%0ANHA"
	msg = strings.Replace(msg, "ABCDEF", otp, -1)
	username := SMS_Username
	password := SMS_Password
	entityID := "1001548700000010184"
	tempID := "1007170748130898041"
	source := "NHASMS"

	urlStr := fmt.Sprintf("https://sms6.rmlconnect.net/bulksms/bulksms?username=%s&password=%s&type=0&dlr=1&destination=%s&source=%s&message=%s&entityid=%s&tempid=%s",
		username, password, phoneNo, source, msg, entityID, tempID)

	response, err := http.Post(urlStr, "application/json", nil)
	if err != nil {
		ErrorLogger.Printf("SMS API failed with error : ")
		ErrorLogger.Println(err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return true
	}

	ErrorLogger.Println("Failed to send SMS. Status code : ", response.StatusCode)
	return false
}


func sendSMSAPInew(phoneNo, otp string) bool {
	if otp=="123456"{
		return true
	}
	msg := "Dear%20User%2C%0AYour%20OTP%20to%20access%20CAPF%20application%20is%20ABCDEF.%20It%20will%20be%20valid%20for%203%20minutes.%0ANHA"
	msg = strings.Replace(msg, "ABCDEF", otp, -1)
	username := SMS_Username
	password := SMS_Password
	entityID := "1001548700000010184"
	tempID := "1007170748130898041"
	source := "NHASMS"
	
	
	payload := map[string]string{
		"userid":   username,
		"password": password,
		"mobile": phoneNo,
		"senderid": source,
		"dltEntityId": entityID,
		"msg": msg,
		"sendMethod": "quick",
		"msgType": "text",
		"dltTemplateId":tempID,
		"output": "json",
		"duplicatecheck": "true",
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return false
	}
	
	urlStr := "http://172.105.57.57/SMSApi/send"

	// client := &http.Client{
	// 	Timeout: 10 * time.Second,
	// }

	response, err := http.Post(urlStr, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		ErrorLogger.Printf("SMS API failed with error : ")
		ErrorLogger.Println(err)
		return false
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return false
	}
	bodyString := string(bodyBytes)

	fmt.Println("Response Body:", bodyString)

	if response.StatusCode == http.StatusOK {
		return true
	}

	ErrorLogger.Println("Failed to send SMS. Status code : ", response.StatusCode)
	return false
}
