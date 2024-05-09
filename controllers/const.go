package controllers

import (
	"log"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	Beneficiary_URL = "https://apis.pmjay.gov.in/prodbis/capfService/v2.0/searchFamilyDetails"
	SMS_Username = getEnv("SMS_USER")
	SMS_Password = getEnv("SMS_PASS")
)