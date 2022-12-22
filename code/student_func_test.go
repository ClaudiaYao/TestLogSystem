package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAdminLogin(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()
	AdminLogin(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if !strings.Contains(string(data), "Password:") {
		t.Error("expected admin page", string(data))
	}
}

func TestLoginGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()
	Login(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if !strings.Contains(string(data), "Please login to your account") {
		t.Errorf("expected ABC got %v", string(data))
	}
}

const (
	FIRSTTIMELOGGING = 0
	ALREADYLOGGING   = 1
	SUBMITTED        = 2
)

func getRandomStudent(logging_status int) *UserInfo {
	if logging_status == FIRSTTIMELOGGING {
		for _, login_record := range LoginMap {
			if login_record.LoggingTime == "" {
				student_id := login_record.StudentID
				return UserMap[student_id]
			}
		}
	} else if logging_status == ALREADYLOGGING {
		for _, login_record := range LoginMap {
			if login_record.LoggingTime != "" && login_record.SubmittingTime == "" {
				student_id := login_record.StudentID
				return UserMap[student_id]
			}
		}

	} else if logging_status == SUBMITTED {
		for _, login_record := range LoginMap {
			if login_record.SubmittingTime != "" {
				student_id := login_record.StudentID
				return UserMap[student_id]
			}
		}
	}
	return &UserInfo{}
}

func TestLoginSubmit(t *testing.T) {
	userInfo := getRandomStudent(ALREADYLOGGING)
	if userInfo == nil {
		t.Errorf("could not find the first logging user")
		return
	}

	// when we create a new test server, we need to specify the
	// handler function. Pay attention that the original mux definition
	// may not work if using the gorilla/mux, therefore, we use
	// the default http.NewServeMux.

	// based on the design of the code, if the redirecting page
	// path should also be defined in the servemux.
	router := http.NewServeMux()
	router.HandleFunc("/login", Login)
	router.HandleFunc("/ExamPage/"+userInfo.StudentID, ExamPage)
	server := httptest.NewServer(router)
	defer server.Close()

	// here, create a new http.Client and set its timeout, instead of
	// using the defautl http.Client
	// netClient := &http.Client{
	// 	Timeout: time.Second * 10,
	// }
	netClient := server.Client()

	// // server.URL is an output, which is appointed by the go program automatically
	form := url.Values{}
	form.Add("UserName", userInfo.UserName)
	form.Add("StudentID", userInfo.StudentID)

	res, err := netClient.PostForm(server.URL+"/login", form)
	// res, err := netClient.PostForm(server.URL+"/login", url.Values{"UserName": {userInfo.UserName}, "StudentID": {userInfo.StudentID}})
	if err != nil {
		log.Fatal("error:", err)
		return
	}

	// here this case will definitely fail because when redirecting to the
	// Exam page, the program checks if the user has been logged in by checking
	// cookie and session. If use postform to test, the cookie might not be
	// generated and redirect to login page again.
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if !strings.Contains(string(data), "Test Main Page") {
		t.Error("expected enter exam page.")
	}
}
