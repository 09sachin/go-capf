package controllers

import (
	"encoding/json"
	"fmt"
	"bytes"
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

	urlStr := Beneficiary_URL
	// Create JSON payload
	payload := map[string]string{
		"id_type":   force_type,
		"id_number": id,
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	search_response, err := http.Post(urlStr, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		ErrorLogger.Printf("Search API failed")
		return
	}

	defer search_response.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(search_response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return
	}

	detailsArray := result["details"].([]interface{})

	var self_data map[string]interface{}
	fmt.Println(self_data["id"])
	var capfData CapfProdNoImageRefresh
	for _, item := range detailsArray {
		// Convert the item to a map[string]interface{}
		detail := item.(map[string]interface{})

		// Check if the "member_type" is "S"
		if detail["member_type"] == "S" {
			self_data = detail
			capfData.MemberNameEng = detail["member_name_eng"].(string)
			capfData.YearOfBirth = detail["year_of_birth"].(string)
			capfData.DOB = detail["dob"].(string)
			capfData.Gender = detail["gender"].(string)
			capfData.InsertionDate = detail["pmjay_id"].(string)
			capfData.MobileNumber = detail["mobile_number"].(string)
			capfData.Id = detail["id_number"].(string)
			capfData.Image = detail["image"].(string)
			break
		}
	}

	var dataList []CapfProdNoImageRefresh
	dataList = append(dataList, capfData)
	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

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


	urlStr := Beneficiary_URL
	// Create JSON payload
	payload := map[string]string{
		"id_type":   force_type,
		"id_number": id,
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	search_response, err := http.Post(urlStr, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		ErrorLogger.Printf("Search API failed")
		return
	}

	defer search_response.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(search_response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return
	}

	detailsArray := result["details"].([]interface{})

	var self_data map[string]interface{}
	fmt.Println(self_data["id"])
	var capfData UserDetail
	for _, item := range detailsArray {
		// Convert the item to a map[string]interface{}
		detail := item.(map[string]interface{})

		// Check if the "member_type" is "S"
		if detail["member_type"] == "S" {
			self_data = detail
			capfData.MemberNameEng = detail["member_name_eng"].(string)
			capfData.DOB = detail["dob"].(string)
			capfData.Gender = detail["gender"].(string)
			capfData.MobileNumber = detail["mobile_number"].(string)
			capfData.PMJAY = detail["pmjay_id"].(string)
			capfData.Id = detail["id_number"].(string)
			capfData.IdType = detail["id_type"].(string)
			capfData.AccountHolder = detail["account_holder_name"].(string)
			capfData.AccountNumber = detail["bank_account_number"].(string)
			capfData.Ifsc =  detail["ifsc_code"].(string)
			capfData.Bank =  detail["bank_name"].(string)
			capfData.SpouseName =  detail["spouse_name_eng"].(string)
			capfData.FatherName =  detail["father_name_eng"].(string)
			capfData.Unit =  detail["unit_name"].(string)
			break
		}
	}
	var dataList []UserDetail
	
	dataList = append(dataList, capfData)

	jsonData, err := json.MarshalIndent(dataList, "", "    ")

	if err != nil {
		fmt.Println(err)
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

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
		data.ClaimAmt = func() string {
			switch {
			case data.PaidAmt.String != "":
				return data.PaidAmt.String
			case data.AppAmt.String != "":
				return data.AppAmt.String
			default:
				return data.SubAmt.String
			}
		}()

		if (data.WorkflowId.String == "171" || data.WorkflowId.String == "172" || data.WorkflowId.String == "173") {
			data.ClaimStatus = "Rejected";
		  } else if (data.WorkflowId.String == "141" ||
		  	data.WorkflowId.String == "142" ||
		  	data.WorkflowId.String == "143") {
			data.ClaimStatus = "Approved";
		  } else if (data.PaidAmt.String!=""){
			data.ClaimStatus =  "Paid";
		  }else{
			data.ClaimStatus = "Pending"
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
