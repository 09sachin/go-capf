package routes

import (
	"github.com/gorilla/mux"
	"github.com/09sachin/go-capf/controllers"
)

func Init() *mux.Router {
	route := mux.NewRouter()

	route.HandleFunc("/send-otp", controllers.SendOtp).Methods("POST")
	route.HandleFunc("/otp-login", controllers.OtpLogin).Methods("POST")
	route.HandleFunc("/dashboard-data", controllers.DashboardData).Methods("GET")
	route.HandleFunc("/user-details", controllers.UserDetails).Methods("GET")
	route.HandleFunc("/hospital-list", controllers.Hospitals).Methods("GET")
	route.HandleFunc("/filter-hospital", controllers.FilterHospital).Methods("GET")
	route.HandleFunc("/queries", controllers.Queries).Methods("GET")
	route.HandleFunc("/track-case", controllers.TrackCases).Methods("GET")
	route.HandleFunc("/claims", controllers.Claims).Methods("GET")

	return route
}
