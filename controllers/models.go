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
	MemberNameEng CustomString
	YearOfBirth   CustomString
	DOB           CustomString
	Gender        CustomString
	InsertionDate CustomString
	MobileNumber  CustomString
	Id            CustomString
}

type Hospital struct {
	HospName        CustomString
	HospContact     CustomString
	HospLatitude    CustomString
	HospLongitude   CustomString
	EmpanelmentType CustomString
}

type UserDetail struct {
	MemberNameEng CustomString
	DOB           CustomString
	Gender        CustomString
	Id            CustomString
	IdType        CustomString
	PMJAY         CustomString
	Unit          CustomString
	AccountHolder CustomString
	Bank          CustomString
	AccountNumber CustomString
	Ifsc          CustomString
	MobileNumber  CustomString
	FatherName    CustomString
	SpouseName    CustomString
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
	Name         CustomString
	CaseNo       CustomString
	ClaimSubDate CustomString
	Status       CustomString
	SubAmt       CustomString
	AppAmt       CustomString
	PaidAmt      CustomString
	HospName     CustomString
	WorkflowId   CustomString
}

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

type RequestBody struct {
	ForceID 	string `json:"force_id"`
	OTP     	string `json:"otp"`
	ForceType   string `json:"force_type"`
}

type PmjayQuery struct{
	PMJAY string
}
