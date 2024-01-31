package controllers

import (
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/09sachin/go-capf/config"
	// "github.com/09sachin/go-capf/models"
	"net/http"
)

type CapfProdNoImageRefresh struct {
	MemberNameEng   string
	YearOfBirth     int // Assuming it's an integer; adjust based on your schema
	DOB             string
	Gender          string
	InsertionDate   string
	MobileNumber    string
	Id              string
}

func DashboardData(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//id := "913228862"

	id := claims.Username


	dashboardQuery := fmt.Sprintf(`SELECT member_name_eng, year_of_birth, dob, gender,
	 insertion_date, mobile_number, id_number 
	 FROM capf.capf_prod_noimage_refresh 
	 WHERE id_number='%s' and relation_name='Self';`, id)

	rows, _ := config.ExecuteQuery(dashboardQuery)
	
	var dataList []CapfProdNoImageRefresh

	for rows.Next() {
		var data CapfProdNoImageRefresh
		err := rows.Scan(&data.MemberNameEng, &data.YearOfBirth, &data.DOB, &data.Gender, &data.InsertionDate, &data.MobileNumber, &data.Id)
		fmt.Println(err)
		dataList = append(dataList, data)	
	}


	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	fmt.Println(err)
	
	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

}


type UserDetail struct {
	MemberNameEng   string
	DOB             string
	Gender          string
	Id              string
	IdType          string
	PMJAY			string
	Unit			string
	AccountHolder   string
	Bank 			string
	AccountNumber   string
	Ifsc            string
	MobileNumber    string
	FatherName      string
	SpouseName      string
}
func UserDetails(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := claims.Username
	user_details_query := fmt.Sprintf(`select member_name_eng, dob, gender, 
	id_number, id_type, pmjay_id, unit_name, account_holder_name, bank_name, bank_account_number, ifsc_code,
	mobile_number, father_name_eng, spouse_name_eng
	from capf.capf_prod_noimage_refresh where id_number='%s' and relation_name='Self';`, id)
	
	rows, _ := config.ExecuteQuery(user_details_query)
	
	var dataList []UserDetail

	for rows.Next() {
		var data UserDetail
		err := rows.Scan(&data.MemberNameEng, &data.DOB, &data.Gender, &data.Id, &data.IdType, &data.PMJAY, &data.Unit, &data.AccountHolder, &data.Bank, &data.AccountNumber, &data.Ifsc, &data.MobileNumber, &data.FatherName, &data.SpouseName)
		fmt.Println(err)
		dataList = append(dataList, data)	
	}


	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	fmt.Println(err)
	
	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}


type Hospital struct{
	HospName   string 
	HospLatitude   string
	HospLongitude  string
	EmpanelmentType   string
}



func Hospitals(w http.ResponseWriter, r *http.Request) {
	query_params := r.URL.Query()
	num := query_params.Get("page")
	size := query_params.Get("size")
	pageSize , err := strconv.Atoi(size)
	page, err := strconv.Atoi(num)
	offset := (page - 1) * pageSize
	empanelment := query_params.Get("empanelment")
	var empanelment_type string
	if empanelment=="PMJAY"{
		empanelment_type = "('PMJAY and CAPF', 'PMJAY', 'PMJAY and CGHS')"
	}else if empanelment=="CAPF"{
		empanelment_type = "('PMJAY and CAPF', 'Only CAPF','PMJAY and CGHS')"
	}else{
		empanelment_type = "('PMJAY and CAPF', 'PMJAY and CGHS')"
	}
	hospital_query := fmt.Sprintf("select empanelment_type, hosp_name, hosp_latitude, hosp_longitude from  hem_t_hosp_info WHERE empanelment_type in %s and active_yn ='Y' and hosp_status ='Approved' LIMIT %d OFFSET %d", empanelment_type, pageSize, offset)
	rows, _ := config.ExecuteQuery(hospital_query)
	
	var dataList []Hospital

	for rows.Next() {
		var data Hospital
		err := rows.Scan(&data.EmpanelmentType, &data.HospName, &data.HospLatitude, &data.HospLongitude)
		fmt.Println(err)
		dataList = append(dataList, data)	
	}

	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	fmt.Println(err)
	
	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}


	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}


type NearestHospital struct {
	HospName   string 
	HospLatitude   string
	HospLongitude  string
	EmpanelmentType   string
}


func FilterHospital(w http.ResponseWriter, r *http.Request) {
	query_params := r.URL.Query()
	radius := query_params.Get("radius")
	latitude := query_params.Get("latitude")
	longitude := query_params.Get("longitude")
	filter_hosp := fmt.Sprintf(` SELECT hosp_name, hosp_latitude, hosp_longitude, empanelment_type
		FROM hem_t_hosp_info
		WHERE 
			CASE WHEN hosp_latitude ~ '^-?\d+(\.\d+)?$' 
				THEN CAST(hosp_latitude AS DOUBLE PRECISION) 
				ELSE NULL 
			END IS NOT NULL
			AND 
			CASE WHEN hosp_longitude ~ '^-?\d+(\.\d+)?$' 
				THEN CAST(hosp_longitude AS DOUBLE PRECISION) 
				ELSE NULL 
			END IS NOT NULL
			AND 6371 * 2 * ASIN(SQRT(
				POWER(SIN(RADIANS(CAST(hosp_latitude AS DOUBLE PRECISION) - CAST(%s AS DOUBLE PRECISION)) / 2), 2) +
				COS(RADIANS(CAST(18.72 AS DOUBLE PRECISION))) * COS(RADIANS(CAST(hosp_latitude AS DOUBLE PRECISION))) *
				POWER(SIN(RADIANS(CAST(hosp_longitude AS DOUBLE PRECISION) - CAST(%s AS DOUBLE PRECISION)) / 2), 2)
			)) <= %s
			AND empanelment_type IN ('PMJAY and CAPF', 'PMJAY', 'Only CAPF', 'PMJAY and CGHS') 
			AND active_yn = 'Y' 
			AND hosp_status = 'Approved' 
		LIMIT 10;`, latitude, longitude, radius)
	
	rows, _ := config.ExecuteQuery(filter_hosp)
	
	var dataList []NearestHospital

	for rows.Next() {
		var data NearestHospital
		err := rows.Scan(&data.HospName, &data.HospLatitude, &data.HospLongitude, &data.EmpanelmentType)
		fmt.Println(err)
		dataList = append(dataList, data)	
	}

	fmt.Println(dataList)


	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	fmt.Println(err)
	fmt.Println(string(jsonData))
	
	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}


	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

type Query struct{
	Remarks 		string
	SubmissionDate  string
	CaseNo  		string
}

func Queries(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := claims.Username
	fmt.Println(id)

	query := fmt.Sprintf(`select  wa.remarks, 
			 reim.claim_sub_dt, reim.case_no
			from capf.tms_t_case_workflow_audit wa 
			join capf.case_dump_capf_reim_pfms reim 
			on reim.patient_no=wa.transaction_id 
			where wa.current_group_id in ('GP603', 'GPSHA', 'GPMD', 'GPBANK') 
			and reim.card_no in (select pmjay_id from capf.capf_prod_noimage_refresh where id_number='%s') and reim.ben_pending='Y' 
			order by wa.crt_dt `, id)
	
	rows, _ := config.ExecuteQuery(query)

	var dataList []Query

	for rows.Next() {
		var data Query
		err := rows.Scan(&data.Remarks, &data.SubmissionDate, &data.CaseNo)
		fmt.Println(err)
		dataList = append(dataList, data)	
	}

	fmt.Println(dataList)


	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	fmt.Println(err)
	fmt.Println(string(jsonData))
	
	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}


	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}


type TrackCase struct{
	CaseNo string
	ClaimSubmissionDate string
	Status string
}
func TrackCases(w http.ResponseWriter, r *http.Request) {
	query_params := r.URL.Query()
	case_no := query_params.Get("case_no")
	track_query := fmt.Sprintf(`select case_no, claim_sub_dt, workflow_status_desc 
	from capf.case_dump_capf_reim_pfms 
	where case_no='%s'`, case_no)
	rows, _ := config.ExecuteQuery(track_query)
	
	var dataList []TrackCase

	for rows.Next() {
		var data TrackCase
		err := rows.Scan(&data.CaseNo, &data.ClaimSubmissionDate, &data.Status)
		fmt.Println(err)
		dataList = append(dataList, data)	
	}

	fmt.Println(dataList)


	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	fmt.Println(err)
	fmt.Println(string(jsonData))
	
	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}


	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}


func UserClaims(w http.ResponseWriter, r *http.Request) {
	query_params := r.URL.Query()
	case_no := query_params.Get("case_no")
	track_query := fmt.Sprintf(`select case_no, claim_sub_dt, workflow_status_desc 
	from capf.case_dump_capf_reim_pfms 
	where case_no='%s'`, case_no)
	fmt.Println((track_query))
	http.Redirect(w, r, "/", 301)
}