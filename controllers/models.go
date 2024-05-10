package controllers

import (
	"database/sql"
	"time"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"os"
)

func getEnv(key string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        ErrorLogger.Printf("Warning: Environment variable %s is not set.\n", key)
        return ""
    }
    return value
}
var secret_key  = getEnv("SECRET_KEY")

var JwtKey = []byte(secret_key)

type CustomString struct {
	sql.NullString
}


type CapfProdNoImageRefresh struct {
	MemberNameEng string
	YearOfBirth   string
	DOB           string
	Gender        string
	InsertionDate string
	MobileNumber  string
	Id            string
	Image         string
}

type Hospital struct {
	HospName        CustomString
	HospContact     CustomString
	HospLatitude    CustomString
	HospLongitude   CustomString
	EmpanelmentType CustomString
}

type UserDetail struct {
	MemberNameEng string
	DOB           string
	Gender        string
	Id            string
	IdType        string
	PMJAY         string
	Unit          string
	AccountHolder string
	Bank          string
	AccountNumber string
	Ifsc          string
	MobileNumber  string
	FatherName    string
	SpouseName    string
}

type NearestHospital struct {
	HospName        CustomString
	HospContact     CustomString
	HospLatitude    CustomString
	HospLongitude   CustomString
	EmpanelmentType CustomString
}


type Query struct {
	Remarks        CustomString
	SubmissionDate CustomString
	CaseNo         CustomString
}

type TrackCase struct {
	CaseNo              CustomString
	ClaimSubmissionDate CustomString
	Status              CustomString
	WorkflowDate        CustomString
}

type UserClaim struct {
	Name         string
	CaseNo       CustomString
	ClaimSubDate CustomString
	Status       CustomString
	SubAmt       CustomString
	AppAmt       CustomString
	PaidAmt      CustomString
	HospName     CustomString
	WorkflowId   CustomString
	ClaimAmt     string
	ClaimStatus  string
}

type TokenClaim struct {
	Username string `json:"username"`
	PmjayId string `json:"pmjayid"`
	Names string `json:"names"`
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

type RequestBody struct {
	ForceID 	string `json:"force_id"`
	OTP     	string `json:"otp"`
	ForceType   string `json:"force_type"`
}

type PmjayQuery struct{
	PMJAY string
}


type ApiResponse struct {
    Details []struct {
        MemberType string `json:"member_type"`
        // Add other fields as needed
    } `json:"details"`
}