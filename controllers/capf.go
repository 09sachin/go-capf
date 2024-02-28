package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/09sachin/go-capf/config"
	"strconv"
	// "github.com/09sachin/go-capf/models"
	"net/http"
)

type CapfProdNoImageRefresh struct {
	MemberNameEng string
	YearOfBirth   int // Assuming it's an integer; adjust based on your schema
	DOB           string
	Gender        string
	InsertionDate string
	MobileNumber  string
	Id            string
}

func DashboardData(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	//id := "913228862"

	id := claims.Username
	force_type := claims.ForceType

	dashboardQuery := fmt.Sprintf(`SELECT member_name_eng, year_of_birth, dob, gender,
	 insertion_date, mobile_number, id_number 
	 FROM capf.capf_prod_noimage_refresh 
	 WHERE id_number='%s' and id_type='%s' and relation_name='Self';`, id, force_type)

	rows, sql_error := config.ExecuteQuery(dashboardQuery)
	if sql_error != nil {
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

func UserDetails(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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
	from capf.capf_prod_noimage_refresh where id_number='%s' and id_type='%s' and relation_name='Self';`, id, force_type)

	rows, sql_error := config.ExecuteQuery(user_details_query)
	if sql_error != nil {
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

type Hospital struct {
	HospName        string
	HospContact		string
	HospLatitude    string
	HospLongitude   string
	EmpanelmentType string
}

func Hospitals(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	_, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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
		empanelment_type = "('PMJAY and CAPF', 'Only CAPF','PMJAY and CGHS')"
	} else {
		empanelment_type = "('PMJAY and CAPF', 'PMJAY and CGHS')"
	}

	var hospital_query string
	if distict!=""{
		hospital_query = fmt.Sprintf("select empanelment_type, hosp_name, hosp_contact_no, hosp_latitude, hosp_longitude from  hem_t_hosp_info WHERE empanelment_type in %s and active_yn ='Y' and hosp_status ='Approved' and state='%s' and district='%s' LIMIT %d OFFSET %d", empanelment_type, state, distict, pageSize, offset)
	}else{
		hospital_query = fmt.Sprintf("select empanelment_type, hosp_name, hosp_contact_no, hosp_latitude, hosp_longitude from  hem_t_hosp_info WHERE empanelment_type in %s and active_yn ='Y' and hosp_status ='Approved' and state='%s' LIMIT %d OFFSET %d", empanelment_type, state, pageSize, offset)
	}

	fmt.Println(hospital_query)

	rows, sql_error := config.ExecuteQuery(hospital_query)
	if sql_error != nil {
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

type NearestHospital struct {
	HospName        string
	HospContact     string
	HospLatitude    string
	HospLongitude   string
	EmpanelmentType string
}

func FilterHospital(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	_, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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
		err := rows.Scan(&data.HospName,&data.HospContact, &data.HospLatitude, &data.HospLongitude, &data.EmpanelmentType)
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

type Query struct {
	Remarks        string
	SubmissionDate string
	CaseNo         string
}

func Queries(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	id := claims.Username
	pmjay := claims.PmjayId
	fmt.Println(id)

	query := fmt.Sprintf(`select  wa.remarks, 
			 reim.claim_sub_dt, reim.case_no
			from capf.tms_t_case_workflow_audit wa 
			join capf.case_dump_capf_reim_pfms reim 
			on reim.patient_no=wa.transaction_id 
			where wa.current_group_id in ('GP603', 'GPSHA', 'GPMD', 'GPBANK') 
			and reim.card_no in %s and reim.ben_pending='Y' 
			order by wa.crt_dt `, pmjay)

	rows, sql_error := config.ExecuteQuery(query)
	if sql_error != nil {
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

type TrackCase struct {
	CaseNo              string
	ClaimSubmissionDate string
	Status              string
	WorkflowDate        string
}

func TrackCases(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	query_params := r.URL.Query()
	case_no := query_params.Get("case_no")
	claims, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	pmjay := claims.PmjayId
	track_query := fmt.Sprintf(`SELECT 
		reimb.case_no,
		reimb.claim_sub_dt,
		workflow.process_desc,
		wa.crt_dt
	FROM 
		capf.tms_t_case_workflow_audit wa
	JOIN 
		capf.case_dump_capf_reim_pfms reimb
	ON 
		wa.transaction_id = reimb.patient_no
	JOIN (
		WITH RankedWorkflows AS (
			SELECT 
				*,
				ROW_NUMBER() OVER (PARTITION BY workflow_id ORDER BY crt_dt) AS row_num
			FROM 
				master.tms_m_case_workflow_new
			WHERE 
				state_code = '91'
		)
		SELECT 
			*
		FROM 
			RankedWorkflows
		WHERE 
			row_num = 1
	) workflow
	ON 
		wa.next_workflow_id = workflow.workflow_id
	WHERE 
		reimb.case_no = '%s' and 
		reimb.card_no in %s
	ORDER BY 
    wa.crt_dt DESC;`, case_no, pmjay)
	rows, sql_error := config.ExecuteQuery(track_query)
	if sql_error != nil {
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

type UserClaim struct {
	Name         string
	CaseNo       string
	ClaimSubDate string
	Status       string
	SubAmt       string
	AppAmt       string
	PaidAmt      string
}

func UserClaims(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	claims, err := getClaimsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{
			Error: "Unauthorised request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	id := claims.Username
	force_type := claims.ForceType
	claims_query := fmt.Sprintf(`SELECT 
    usr.member_name_eng, 
    rem.case_no, 
    rem.claim_sub_dt, 
    rw.process_desc,
    rem.claim_sub_amt, 
    rem.claim_app_amt, 
    rem.claim_paid_amt
FROM 
    capf.case_dump_capf_reim_pfms rem
JOIN 
    capf.capf_prod_noimage_refresh usr ON rem.card_no = usr.pmjay_id
JOIN (
	WITH RankedWorkflows AS (
		SELECT 
			*,
			ROW_NUMBER() OVER (PARTITION BY workflow_id ORDER BY crt_dt) AS row_num
		FROM 
			master.tms_m_case_workflow_new
		WHERE 
			state_code = '91'
	)
	SELECT 
		*
	FROM 
		RankedWorkflows
	WHERE 
		row_num = 1
) rw ON rem.workflow_id = rw.workflow_id
WHERE 
    usr.id_number = '%s' and usr.id_type='%s';`, id, force_type)

	rows, sql_error := config.ExecuteQuery(claims_query)
	if sql_error != nil {
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
		err := rows.Scan(&data.Name, &data.CaseNo, &data.ClaimSubDate, &data.Status, &data.SubAmt, &data.AppAmt, &data.PaidAmt)
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
