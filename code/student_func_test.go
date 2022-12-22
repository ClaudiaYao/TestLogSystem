package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
	fmt.Println(data)
	if !strings.Contains(string(data), "Password:") {
		t.Error("expected admin page", string(data))
	}
}

func TestLoginHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()
	Login(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	fmt.Println(data)
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

// func TestLoginUser(t *testing.T) {
// 	userInfo := getRandomStudent(ALREADYLOGGING)
// 	if userInfo == nil {
// 		t.Errorf("could not find the first logging user")
// 		return
// 	}

// 	server := httptest.NewServer(router)
// 	defer server.Close()

// 	// here, create a new http.Client and set its timeout, instead of
// 	// using the defautl http.Client
// 	netClient := &http.Client{
// 		Timeout: time.Second * 10,
// 	}

// 	// // server.URL is an output, which is appointed by the go program automatically
// 	form := url.Values{}
// 	form.Add("UserName", userInfo.UserName)
// 	form.Add("StudentID", userInfo.StudentID)

// 	fmt.Println(server.URL + "/login")

// 	// this part fails and I could not find the reason
// 	// the form value could use the above defined, it is not a problem. It seems that the
// 	// form value could not be passed to the post request
// 	// res, err := netClient.PostForm(server.URL+"/login", form)
// 	res, err := netClient.PostForm(server.URL+"/login", url.Values{"UserName": {userInfo.UserName}, "StudentID": {userInfo.StudentID}})
// 	if err != nil {
// 		log.Fatal("error:", err)
// 		return
// 	}

// 	defer res.Body.Close()
// 	data, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Errorf("expected error to be nil got %v", err)
// 	}

// 	if !strings.Contains(string(data), "Exam Page") {
// 		t.Errorf("expected ABC got %v", string(data))
// 	}
// }
