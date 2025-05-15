package controllers

import (
	"log"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	Beneficiary_URL = "https://apisprod.nha.gov.in/pmjay/prodbis/capfService/v3.0/searchFamilyDetails"
	SMS_Username = "PMJCAPF"
	SMS_Password = "pybnZ4tD"
	CLAIMS_UPDATE_BASE_URL = getEnv("CLAIMS_UPDATE_BASE_URL")
	PMJAY_CLIENT_ID = getEnv("PMJAY_CLIENT_ID")
	PMJAY_CLIENT_SECRET = getEnv("PMJAY_CLIENT_SECRET")
	PMJAY_CLIENT_USERNAME = getEnv("PMJAY_CLIENT_USERNAME")
	PMJAY_CLIENT_PASSWORD = getEnv("PMJAY_CLIENT_PASSWORD")
)
