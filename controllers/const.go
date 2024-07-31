package controllers

import (
	"log"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	Beneficiary_URL = getEnv("SEARCH_URL")
	SMS_Username = getEnv("SMS_USER")
	SMS_Password = getEnv("SMS_PASS")
	CLAIMS_UPDATE_BASE_URL = getEnv("CLAIMS_UPDATE_BASE_URL")
)