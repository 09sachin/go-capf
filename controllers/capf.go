package controllers

import (
	"fmt"
	"encoding/json"
	// "github.com/09sachin/go-capf/config"
	// "github.com/09sachin/go-capf/models"
	"net/http"
)


func DashboardData(w http.ResponseWriter, r *http.Request) {

	dashboard_query := "select member_name_eng, year_of_birth, dob, gender, insertion_date, mobile_number from capf.capf_prod_noimage_refresh  where id_number='000000523';"
	fmt.Println(dashboard_query)
	response := Response{
		Message: "Hello, JSON!",
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response writer
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

}


func UserDetails(w http.ResponseWriter, r *http.Request) {

	user_details_query := "select member_name_eng, year_of_birth,dob, gender, insertion_date, pfms_id, mobile_number,address, pincode, state_lgd_code, district_lgd_code, subdistrict_lgd_code, village_town_lgd_code from capf.capf_prod_noimage_refresh limit 10"
	fmt.Println(user_details_query)

	http.Redirect(w, r, "/", 301)
}



func Hospitals(w http.ResponseWriter, r *http.Request) {

	hospital_query := "select * from  hem_t_hosp_info WHERE empanelment_type in ( 'PMJAY and CAPF', 'PMJAY','Only CAPF','PMJAY and CGHS') and active_yn ='Y' and hosp_status ='Approved'"
	fmt.Println(hospital_query)

	http.Redirect(w, r, "/", 301)
}

func FilterHospital(w http.ResponseWriter, r *http.Request) {
	filter_hosp := ` SELECT hosp_name, hosp_latitude, hosp_longitude
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
				POWER(SIN(RADIANS(CAST(hosp_latitude AS DOUBLE PRECISION) - CAST(18.72 AS DOUBLE PRECISION)) / 2), 2) +
				COS(RADIANS(CAST(18.72 AS DOUBLE PRECISION))) * COS(RADIANS(CAST(hosp_latitude AS DOUBLE PRECISION))) *
				POWER(SIN(RADIANS(CAST(hosp_longitude AS DOUBLE PRECISION) - CAST(79.97 AS DOUBLE PRECISION)) / 2), 2)
			)) <= 10 
			AND empanelment_type IN ('PMJAY and CAPF', 'PMJAY', 'Only CAPF', 'PMJAY and CGHS') 
			AND active_yn = 'Y' 
			AND hosp_status = 'Approved' 
		LIMIT 10;`
	fmt.Println(filter_hosp)

	http.Redirect(w, r, "/", 301)
}


func Queries(w http.ResponseWriter, r *http.Request) {

	query := "select wa.transaction_id, wa.remarks, wa.current_group_id, wa.crt_dt, reim.claim_sub_dt from capf.tms_t_case_workflow_audit wa join capf.case_dump_capf_reim_pfms reim on reim.patient_no=wa.transaction_id where wa.current_group_id in ('GP603', 'GPSHA', 'GPMD', 'GPBANK') and reim.card_no='PG1OU04V2' and reim.ben_pending='Y' order by wa.crt_dt limit 1"
	fmt.Println(query)

	http.Redirect(w, r, "/", 301)
}


func TrackCases(w http.ResponseWriter, r *http.Request) {

	track_query := "select case_no, claim_sub_dt, workflow_status_desc from capf.case_dump_capf_reim_pfms where case_no='REM/2022/470381'"
	fmt.Println(track_query)

	http.Redirect(w, r, "/", 301)
}

func Claims(w http.ResponseWriter, r *http.Request) {

	fmt.Println("opt sent")

	http.Redirect(w, r, "/", 301)
}