package controllers

import (
	"log"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	Beneficiary_URL = "https://apis-uat.pmjay.gov.in/pmjay/bis/capfService/searchFamilyDetails"
	SMS_Username = getEnv("SMS_USER")
	SMS_Password = getEnv("SMS_PASS")
)