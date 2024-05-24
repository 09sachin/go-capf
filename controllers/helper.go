package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"net/http"
)

func (cs CustomString) MarshalJSON() ([]byte, error) {
	if cs.Valid {
		return json.Marshal(cs.String)
	}
	return json.Marshal("") // Convert null to empty string in JSON
}


func (cs *CustomString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	cs.Valid = true
	cs.String = str
	return nil
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

func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}


func ParseInt (s string) (int, error) {
	return strconv.Atoi(s)
}

func isAlphaNumeric(s string) bool {
	for _, char := range s {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >='0' && char <='9') || char=='/' || char=='-') {
			return false
		}
	}
	return true
}

func QueryParamsError(w http.ResponseWriter){
	w.WriteHeader(http.StatusUnprocessableEntity)
	response := ErrorResponse{
		Error: "Unexpected value in query params",
	}
	json.NewEncoder(w).Encode(response)
}


func UnauthorisedError(w http.ResponseWriter){
	w.WriteHeader(http.StatusUnauthorized)
	response := ErrorResponse{
		Error: "Unauthorised request",
	}
	json.NewEncoder(w).Encode(response)
}


func JsonParseError(w http.ResponseWriter){
	w.WriteHeader(http.StatusNotFound)
	response := ErrorResponse{
		Error: "Json Parse error",
	}
	json.NewEncoder(w).Encode(response)
}


func JsonEncodeError(w http.ResponseWriter){
	w.WriteHeader(http.StatusNotFound)
	response := ErrorResponse{
		Error: "Encoder error",
	}
	json.NewEncoder(w).Encode(response)
}


func DbError(w http.ResponseWriter){
	w.WriteHeader(http.StatusNotFound)
	response := ErrorResponse{
		Error: "Internal server error",
	}
	json.NewEncoder(w).Encode(response)
}

//Custom 404 error
func Custom4O4Error(w http.ResponseWriter, s string){
	w.WriteHeader(http.StatusNotFound)
	response := ErrorResponse{
		Error: s,
	}
	json.NewEncoder(w).Encode(response)
}


func maskPhoneNumber(phone string) string {
    // Check if the phone number length is less than 10
    if len(phone) < 10 {
        return "Invalid phone number"
    }

    // Get the last four characters of the phone number
    lastFour := phone[len(phone)-4:]

    // Create a masked string with 'x' for the first six characters
    masked := "XXXXXX" + lastFour

    return masked
}