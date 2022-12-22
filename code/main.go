package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"

	// "math/rand"
	"net/http"
	// "github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

var tpl *template.Template
var mapSessions = map[string]string{}

// var sessions = map[string]session{}
var errInvalidLogInError = errors.New("InvalidLogInError: username or password does not match")

const COOKIE_NAME = "Test_Cookie"

var UserMap = make(map[string]*UserInfo)
var LoginMap = make(map[string]*LoginInfo)
var Admin AdminInfo
var router *mux.Router

// generate_sessions() function is only used to generate sessions for testing purpose
// func generate_sessions() {
// 	users := []string{}
// 	for k, _ := range UserMap {
// 		if rand.Intn(5) > 2 {
// 			users = append(users, k)
// 		}
// 	}
// 	user_count := len(users)
// 	for i := 0; i < user_count; i++ {
// 		session, _ := uuid.NewV4()
// 		session_id := session.String()
// 		mapSessions[session_id] = users[i]

// 	}

// }

// the first thing for the server do is to load the User, Appointment and Venue information from those JSON files.
func init() {
	// SaveToTXTFiles()

	//This function only runs when you want to change the password of the admin.
	// For this demo project, use Admin UserName: Claudia, Admin Pwd: Admin
	// generate_admin_data()

	tpl = ParseTemplates()

	// generate_sessions()
	ReadFromUserFile()
	ReadFromLoginFile()
	ReadJsonAdmin()
	fmt.Println(Admin)
	createSubmissionFolder()

}

func main() {

	router = mux.NewRouter()
	// files := http.FileServer(http.Dir(config.Static))
	// mux.Handle("/static/", http.StripPrefix("/static/", files))

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	router.HandleFunc("/", Index)
	router.HandleFunc("/admin", AdminLogin)
	router.HandleFunc("/login", Login)
	router.HandleFunc("/ExamPage/{StudentID}", ExamPage).Methods("GET", "PUT", "POST", "DELETE")
	router.HandleFunc("/submitted", Submitted)
	router.HandleFunc("/admin/setting", AdminSetting)
	router.HandleFunc("/admin/upload", AdminUpload)
	// router.HandleFunc("/api/v1/courses/{courseid}", course).Methods(
	// 	"GET", "PUT", "POST", "DELETE")
	fmt.Println("Listening at port 8080")
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(server.ListenAndServe())

	// http.ListenAndServeTLS("localhost:8080", "./cert/cert.pem", "./cert/key.pem", nil)
}

// the handler function when visiting Index page.
func Index(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/login", http.StatusSeeOther)
}
