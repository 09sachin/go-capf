package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/09sachin/go-capf/config"
	"net/http"
)


func DashboardData(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	id := claims.Username
	force_type := claims.ForceType

	dashboardQuery := fmt.Sprintf(`SELECT member_name_eng, year_of_birth, dob, gender,
	 insertion_date, mobile_number, id_number 
	 FROM capf_prod_noimage_refresh 
	 WHERE id_number='%s' and id_type='%s' and relation_name='Self';`, id, force_type)

	rows, sql_error := config.ExecuteQuery(dashboardQuery)
	if sql_error != nil {
		ErrorLogger.Printf("Database connection error : dashboard")
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var dataList []CapfProdNoImageRefresh

	for rows.Next() {
		var data CapfProdNoImageRefresh
		err := rows.Scan(&data.MemberNameEng, &data.YearOfBirth, &data.DOB, &data.Gender, &data.InsertionDate, &data.MobileNumber, &data.Id)
		if err != nil {
			fmt.Println(err)
		}
		dataList = append(dataList, data)
	}

	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error encoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

}


func UserDetails(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	id := claims.Username
	force_type := claims.ForceType
	user_details_query := fmt.Sprintf(`select member_name_eng, dob, gender, 
	id_number, id_type, pmjay_id, unit_name, account_holder_name, bank_name, bank_account_number, ifsc_code,
	mobile_number, father_name_eng, spouse_name_eng
	from capf_prod_noimage_refresh where id_number='%s' and id_type='%s' and relation_name='Self';`, id, force_type)

	rows, sql_error := config.ExecuteQuery(user_details_query)
	if sql_error != nil {
		ErrorLogger.Printf("Database connection error : userdetails")
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var dataList []UserDetail

	for rows.Next() {
		var data UserDetail
		err := rows.Scan(&data.MemberNameEng, &data.DOB, &data.Gender, &data.Id, &data.IdType, &data.PMJAY, &data.Unit, &data.AccountHolder, &data.Bank, &data.AccountNumber, &data.Ifsc, &data.MobileNumber, &data.FatherName, &data.SpouseName)
		if err != nil {
			fmt.Println(err)
		}
		dataList = append(dataList, data)
	}

	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error encoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}


func Hospitals(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	_, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	query_params := r.URL.Query()
	num := query_params.Get("page")
	size := query_params.Get("size")
	state := query_params.Get("state")
	distict := query_params.Get("district")
	pageSize, _ := strconv.Atoi(size)
	page, _ := strconv.Atoi(num)
	offset := (page - 1) * pageSize
	empanelment := query_params.Get("empanelment")
	var empanelment_type string
	if empanelment == "PMJAY" {
		empanelment_type = "('PMJAY and CAPF', 'PMJAY', 'PMJAY and CGHS')"
	} else if empanelment == "CAPF" {
		empanelment_type = "('PMJAY and CAPF', 'Only CAPF','PMJAY and CGHS', 'Only CGHS')"
	} else {
		empanelment_type = "('PMJAY and CAPF', 'PMJAY and CGHS')"
	}

	var hospital_query string
	if distict != "" {
		hospital_query = fmt.Sprintf("select empanelment_type, hosp_name, hosp_contact_no, hosp_latitude, hosp_longitude from  hem_t_hosp_info WHERE empanelment_type in %s and active_yn ='Y' and hosp_status ='Approved' and state='%s' and district='%s' LIMIT %d OFFSET %d", empanelment_type, state, distict, pageSize, offset)
	} else {
		hospital_query = fmt.Sprintf("select empanelment_type, hosp_name, hosp_contact_no, hosp_latitude, hosp_longitude from  hem_t_hosp_info WHERE empanelment_type in %s and active_yn ='Y' and hosp_status ='Approved' and state='%s' LIMIT %d OFFSET %d", empanelment_type, state, pageSize, offset)
	}

	rows, sql_error := config.ExecuteQuery(hospital_query)
	if sql_error != nil {
		ErrorLogger.Printf("Database connection error : hospitals")
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var dataList []Hospital

	for rows.Next() {
		var data Hospital
		err := rows.Scan(&data.EmpanelmentType, &data.HospName, &data.HospContact, &data.HospLatitude, &data.HospLongitude)
		if err != nil {
			fmt.Println(err)
		}
		dataList = append(dataList, data)
	}

	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error encoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}


func FilterHospital(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	_, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	query_params := r.URL.Query()
	radius := query_params.Get("radius")
	latitude := query_params.Get("latitude")
	longitude := query_params.Get("longitude")
	filter_hosp := fmt.Sprintf(` SELECT hosp_name, hosp_contact_no, hosp_latitude, hosp_longitude, empanelment_type
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

	rows, err := config.ExecuteQuery(filter_hosp)
	if err != nil {
		ErrorLogger.Printf("Database connection error : filter hospital")
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	var dataList []NearestHospital

	for rows.Next() {
		var data NearestHospital
		err := rows.Scan(&data.HospName, &data.HospContact, &data.HospLatitude, &data.HospLongitude, &data.EmpanelmentType)
		if err != nil {
			fmt.Println(err)
			continue
		}
		dataList = append(dataList, data)
	}

	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error encoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}


func Queries(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	id := claims.Username
	pmjay := claims.PmjayId
	fmt.Println(id)

	query := fmt.Sprintf(`select remarks, 
			claim_sub_dt, case_no
			from queries 
			where card_no in %s
			order by crt_dt `, pmjay)

	rows, sql_error := config.ExecuteQuery(query)
	if sql_error != nil {
		ErrorLogger.Printf("Database connection error : Queries")
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var dataList []Query

	for rows.Next() {
		var data Query
		err := rows.Scan(&data.Remarks, &data.SubmissionDate, &data.CaseNo)
		if err != nil {
			fmt.Println(err)
		}
		dataList = append(dataList, data)
	}

	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error encoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}


func TrackCases(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	query_params := r.URL.Query()
	case_no := query_params.Get("case_no")
	claims, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	pmjay := claims.PmjayId
	track_query := fmt.Sprintf(`SELECT 
		case_no,
		claim_sub_dt,
		process_desc,
		crt_dt from 
	track_case
	WHERE 
		case_no = '%s' and 
		card_no in %s
	ORDER BY 
    crt_dt DESC;`, case_no, pmjay)
	rows, sql_error := config.ExecuteQuery(track_query)
	if sql_error != nil {
		ErrorLogger.Printf("Database connection error : trackcase")
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var dataList []TrackCase
	for rows.Next() {
		var data TrackCase
		err := rows.Scan(&data.CaseNo, &data.ClaimSubmissionDate, &data.Status, &data.WorkflowDate)
		if err != nil {
			fmt.Println(err)
		}
		dataList = append(dataList, data)
	}

	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error encoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}



func UserClaims(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	id := claims.Username
	force_type := claims.ForceType
	claims_query := fmt.Sprintf(`select distinct
    member_name_eng, 
    case_no, 
    claim_sub_dt, 
    process_desc,
    claim_sub_amt, 
    claim_app_amt, 
    claim_paid_amt,
	workflow_id,
	hosp_name
FROM 
    claims
WHERE 
    id_number = '%s' and id_type='%s';`, id, force_type)

	rows, sql_error := config.ExecuteQuery(claims_query)
	if sql_error != nil {
		ErrorLogger.Printf("Database connection error : userclaims")
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Database connection could not be established",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	var dataList []UserClaim

	for rows.Next() {
		var data UserClaim
		err := rows.Scan(&data.Name, &data.CaseNo, &data.ClaimSubDate, &data.Status, &data.SubAmt, &data.AppAmt, &data.PaidAmt, &data.WorkflowId, &data.HospName)
		if err != nil {
			fmt.Println(err)
		}
		dataList = append(dataList, data)
	}

	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	// Encode the response as JSON and write it to the response writer
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Error encoding JSON",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}
