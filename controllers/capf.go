package controllers

import (
	"fmt"
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
}

func DashboardData(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := claims.Username


	dashboardQuery := fmt.Sprintf(`SELECT member_name_eng, year_of_birth, dob, gender,
	 insertion_date, mobile_number 
	 FROM capf.capf_prod_noimage_refresh 
	 WHERE id_number='%s' and relation_name='Self';`, id)

	rows, _ := config.ExecuteQuery(dashboardQuery)
	
	var dataList []CapfProdNoImageRefresh

	for rows.Next() {
		var data CapfProdNoImageRefresh
		err := rows.Scan(&data.MemberNameEng, &data.YearOfBirth, &data.DOB, &data.Gender, &data.InsertionDate, &data.MobileNumber)
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


func UserDetails(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := claims.Username
	user_details_query := fmt.Sprintf(`select member_name_eng, year_of_birth,dob, gender, 
	insertion_date, pfms_id, mobile_number,address, pincode, state_lgd_code, district_lgd_code, 
	subdistrict_lgd_code, village_town_lgd_code 
	from capf.capf_prod_noimage_refresh where id_number='%s' and relation_name='Self';`, id)
	fmt.Println(user_details_query)

	http.Redirect(w, r, "/", 301)
}



func Hospitals(w http.ResponseWriter, r *http.Request) {

	hospital_query := "select * from  hem_t_hosp_info WHERE empanelment_type in ( 'PMJAY and CAPF', 'PMJAY','Only CAPF','PMJAY and CGHS') and active_yn ='Y' and hosp_status ='Approved'"
	fmt.Println(hospital_query)

	http.Redirect(w, r, "/", 301)
}

func FilterHospital(w http.ResponseWriter, r *http.Request) {
	query_params := r.URL.Query()
	radius := query_params.Get("radius")
	filter_hosp := fmt.Sprintf(` SELECT hosp_name, hosp_latitude, hosp_longitude
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
			)) <= %s
			AND empanelment_type IN ('PMJAY and CAPF', 'PMJAY', 'Only CAPF', 'PMJAY and CGHS') 
			AND active_yn = 'Y' 
			AND hosp_status = 'Approved' 
		LIMIT 10;`, radius)
	fmt.Println(filter_hosp)

	http.Redirect(w, r, "/", 301)
}


func Queries(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := claims.Username

	query := fmt.Sprintf(`select wa.transaction_id, wa.remarks, 
			wa.current_group_id, wa.crt_dt, reim.claim_sub_dt 
			from capf.tms_t_case_workflow_audit wa 
			join capf.case_dump_capf_reim_pfms reim 
			on reim.patient_no=wa.transaction_id 
			where wa.current_group_id in ('GP603', 'GPSHA', 'GPMD', 'GPBANK') 
			and reim.card_no='%s' and reim.ben_pending='Y' 
			order by wa.crt_dt limit 1`, id)
	fmt.Println(query)

	http.Redirect(w, r, "/", 301)
}


func TrackCases(w http.ResponseWriter, r *http.Request) {
	query_params := r.URL.Query()
	case_no := query_params.Get("case_no")
	track_query := fmt.Sprintf(`select case_no, claim_sub_dt, workflow_status_desc 
	from capf.case_dump_capf_reim_pfms 
	where case_no='%s'`, case_no)
	fmt.Println(track_query)
	http.Redirect(w, r, "/", 301)
}

func UserClaims(w http.ResponseWriter, r *http.Request) {

	fmt.Println("opt sent")

	http.Redirect(w, r, "/", 301)
}