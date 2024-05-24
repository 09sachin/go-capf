package controllers

import (
	"encoding/json"
	"fmt"
	"bytes"
	"github.com/09sachin/go-capf/config"
	"net/http"
	"strings"
)


func DashboardData(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		UnauthorisedError(w)
		return
	}
	
	id := claims.Username
	force_type := claims.ForceType

	//Test conditon for playstore/appstore
	if (id=="00000000" && force_type=="BS"){
		var capfData CapfProdNoImageRefresh
		var dataList []CapfProdNoImageRefresh
		capfData.MemberNameEng = "member_name_eng"
		capfData.YearOfBirth = "year_of_birth"
		capfData.DOB = "01/01/1899"
		capfData.Gender = "gender"
		capfData.InsertionDate = "pmjay_id"
		capfData.MobileNumber = "mobile_number"
		capfData.Id = "id_number"
		capfData.Image = "image"
		dataList = append(dataList, capfData)
		jsonData, err := json.MarshalIndent(dataList, "", "    ")
	
		if err != nil {
			fmt.Println(err)
			JsonParseError(w)
			return
		}
	
		response := JsonResponse{
			Message: json.RawMessage(jsonData),
		}
	
		json.NewEncoder(w).Encode(response)
		return
	}

	urlStr := Beneficiary_URL
	// Create JSON payload
	payload := map[string]string{
		"id_type":   force_type,
		"id_number": id,
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		JsonParseError(w)
		return
	}
	search_response, err := http.Post(urlStr, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		ErrorLogger.Printf("Search API failed")
		Custom4O4Error(w,"Search API failed")
		return
	}

	defer search_response.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(search_response.Body).Decode(&result)
	if err != nil {
		JsonEncodeError(w)
		return
	}

	detailsArray, ok := result["details"].([]interface{})
	if !ok {
		// Handle the case where "details" key is missing
		Custom4O4Error(w,"Details missing in CAPF data")
		return
	}

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
		JsonParseError(w)
		return
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		JsonEncodeError(w)
		return
	}

}


func UserDetails(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		UnauthorisedError(w)
		return
	}

	id := claims.Username
	force_type := claims.ForceType

	//Test conditon for playstore/appstore
	if (id=="00000000" && force_type=="BS"){
		var capfData UserDetail
		var dataList []UserDetail
		capfData.MemberNameEng = "member_name_eng"
		capfData.DOB = "dob"
		capfData.Gender = "gender"
		capfData.MobileNumber = "mobile_number"
		capfData.PMJAY = "pmjay_id"
		capfData.Id = "id_number"
		capfData.IdType = "id_type"
		capfData.AccountHolder = "account_holder_name"
		capfData.AccountNumber = "bank_account_number"
		capfData.Ifsc =  "ifsc_code"
		capfData.Bank =  "bank_name"
		capfData.SpouseName =  "spouse_name_eng"
		capfData.FatherName =  "father_name_eng"
		capfData.Unit =  "unit_name"
		dataList = append(dataList, capfData)
		jsonData, err := json.MarshalIndent(dataList, "", "    ")
	
		if err != nil {
			fmt.Println(err)
			JsonParseError(w)
			return
		}
	
		response := JsonResponse{
			Message: json.RawMessage(jsonData),
		}
	
		json.NewEncoder(w).Encode(response)
		return
	}


	urlStr := Beneficiary_URL
	// Create JSON payload
	payload := map[string]string{
		"id_type":   force_type,
		"id_number": id,
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		JsonParseError(w)
		return
	}
	search_response, err := http.Post(urlStr, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		ErrorLogger.Printf("Search API failed")
		Custom4O4Error(w,"Search API failed")
		return
	}

	defer search_response.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(search_response.Body).Decode(&result)
	if err != nil {
		JsonEncodeError(w)
		return
	}

	detailsArray, ok := result["details"].([]interface{})
	if !ok {
		// Handle the case where "details" key is missing
		Custom4O4Error(w,"Details missing in CAPF data")
		return
	}

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
		JsonParseError(w)
		return
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		JsonEncodeError(w)
		return
	}
}


func Hospitals(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	_, err := getClaimsFromRequest(r)
	if err != nil {
		UnauthorisedError(w)
		return
	}

	query_params := r.URL.Query()
	num := query_params.Get("page")
	size := query_params.Get("size")
	state := query_params.Get("state")
	distict := query_params.Get("district")
	empanelment := query_params.Get("empanelment")
	
	nonAlphanumeric := (isAlphaNumeric(state) && isAlphaNumeric(distict) && isAlphaNumeric(empanelment))
	
	if !nonAlphanumeric{
		QueryParamsError(w)
		return
	}

	pageSize, err := ParseInt(size)
	if err != nil {
		QueryParamsError(w)
		return
	}

	page, err := ParseInt(num)
	if err != nil {
		QueryParamsError(w)
		return
	}

	offset := (page - 1) * pageSize

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
		DbError(w)
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
		JsonParseError(w)
		return
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		JsonEncodeError(w)
		return
	}
}


func FilterHospital(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	_, err := getClaimsFromRequest(r)
	if err != nil {
		UnauthorisedError(w)
		return 
	}

	query_params := r.URL.Query()
	radiusStr := query_params.Get("radius")
	latitudeStr := query_params.Get("latitude")
	longitudeStr := query_params.Get("longitude")

	radius, err := ParseFloat(radiusStr)
	if err != nil {
		QueryParamsError(w)
		return
	}

	latitude, err := ParseFloat(latitudeStr)
	if err != nil {
		QueryParamsError(w)
		return
	}

	longitude, err := ParseFloat(longitudeStr)
	if err != nil {
		QueryParamsError(w)
		return
	}

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
				POWER(SIN(RADIANS(CAST(hosp_latitude AS DOUBLE PRECISION) - CAST(%f AS DOUBLE PRECISION)) / 2), 2) +
				COS(RADIANS(CAST(18.72 AS DOUBLE PRECISION))) * COS(RADIANS(CAST(hosp_latitude AS DOUBLE PRECISION))) *
				POWER(SIN(RADIANS(CAST(hosp_longitude AS DOUBLE PRECISION) - CAST(%f AS DOUBLE PRECISION)) / 2), 2)
			)) <= %f
			AND empanelment_type IN ('PMJAY and CAPF', 'PMJAY', 'Only CAPF', 'PMJAY and CGHS') 
			AND active_yn = 'Y' 
			AND hosp_status = 'Approved' 
		LIMIT 10;`, latitude, longitude, radius)

	rows, err := config.ExecuteQuery(filter_hosp)
	if err != nil {
		DbError(w)
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
		JsonParseError(w)
		return
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		JsonEncodeError(w)
		return
	}
}


func Queries(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		UnauthorisedError(w)
		return
	}

	id := claims.Username
	pmjay := claims.PmjayId
	fmt.Println(id)

	query := `select remarks, 
			claim_sub_dt, case_no
			from queries 
			where card_no in $1
			order by crt_dt `

	rows, sql_error := config.ExecuteQuery(query, pmjay)
	if sql_error != nil {
		DbError(w)
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
		JsonParseError(w)
		return
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		JsonEncodeError(w)
		return
	}
}


func TrackCases(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	query_params := r.URL.Query()
	case_no := query_params.Get("case_no")
	
	isvalid := isAlphaNumeric(case_no)

	if !isvalid{
		QueryParamsError(w)
		return
	}

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		UnauthorisedError(w)
		return
	}

	pmjay := claims.PmjayId
	track_query := `SELECT 
		case_no,
		claim_sub_dt,
		process_desc,
		crt_dt from 
	track_case
	WHERE 
		case_no = $1 and 
		card_no in $2
	ORDER BY 
    crt_dt DESC;`
	rows, sql_error := config.ExecuteQuery(track_query,case_no, pmjay)
	if sql_error != nil {
		DbError(w)
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
		DbError(w)
		return
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		JsonEncodeError(w)
		return
	}
}



func UserClaims(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		UnauthorisedError(w)
		return
	}

	pmjay := claims.PmjayId
	names := claims.Names
	var names_list []string
	elements := strings.Split(names, ", ")
	names_list = append(names_list, elements...)
	var card_list []string
	len_ids := len(pmjay)
	if (len_ids< 2){
		ErrorLogger.Println(pmjay)
		JsonParseError(w)
		return
	}

	pmjay_card := pmjay[1:len_ids-1]
	elements_card := strings.Split(pmjay_card, ", ")
	card_list = append(card_list, elements_card...)
	nameMap := make(map[string]string)

    for i := 0; i < len(names_list); i++ {
		card_len := len(card_list[i])
		names_len := len(names_list[i])
		name_person := names_list[i][1:names_len-1]
		card_no := card_list[i][1:card_len-1]
    	nameMap[card_no] = name_person
    }

	pmjay_card_list := card_list
	InfoLogger.Println(pmjay_card_list)
	placeholders := make([]string, len(pmjay_card_list))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	placeholderStr := strings.Join(placeholders, ", ")
	InfoLogger.Println(placeholderStr)

	claims_query := fmt.Sprintf(`select distinct
    case_no, 
    claim_sub_dt, 
    process_desc,
    claim_sub_amt, 
    claim_app_amt, 
    claim_paid_amt,
	workflow_id,
	hosp_name,
	card_no
	FROM 
		claims
	WHERE 
		card_no in (%s);`, placeholderStr)

	args := make([]interface{}, len(pmjay))
	for i, v := range pmjay {
		args[i] = v
	}
	rows, sql_error := config.ExecuteQuery(claims_query, args...)
	if sql_error != nil {
		ErrorLogger.Println(sql_error)
		DbError(w)
		return
	}
	var dataList []UserClaim

	for rows.Next() {
		var data UserClaim
		err := rows.Scan(&data.CaseNo, &data.ClaimSubDate, &data.Status, &data.SubAmt, &data.AppAmt, &data.PaidAmt, &data.WorkflowId, &data.HospName, &data.CardNo)
		if err != nil {
			fmt.Println(err)
		}
		card_no := data.CardNo.String
		mem_name := nameMap[card_no]
		data.Name = mem_name
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
		JsonParseError(w)
		return
	}

	response := JsonResponse{
		Message: json.RawMessage(jsonData),
	}

	
	errr := json.NewEncoder(w).Encode(response)
	if errr != nil {
		JsonEncodeError(w)
	}
}
