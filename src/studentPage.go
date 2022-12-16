package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/kennygrant/sanitize"
)

func ExamPage(res http.ResponseWriter, req *http.Request) {
	// http.
	// 	fmt.Println(req.URL)
	params := mux.Vars(req)

	// POST is for creating new course
	if req.Method == http.MethodGet {
		// read the string sent to the service
		// that means the client side posts content in JSON format
		student_id := sanitize.HTML(params["StudentID"])
		_, ok := UserMap[student_id]
		if !ok {
			http.Error(res, "Username or student id does not match.", http.StatusUnauthorized)
			fmt.Println("not match!")
			return

		} else {
			user := UserMap[student_id]
			user_login := LoginMap[student_id]
			fmt.Println(user_login)
			start_time, error := time.Parse("2006-01-02 15:04:05", user_login.LoggingTime)

			if error != nil {
				fmt.Println(error)
				return
			}
			end_time := start_time.Add(time.Hour * 2)
			user_info := map[string]interface{}{"UserName": user.UserName, "StudentID": user.StudentID,
				"StartTime": user_login.LoggingTime, "EndTime": end_time}
			tpl.ExecuteTemplate(res, "examPage.html", user_info)
		}

	}
}

func CheckUserInfoMatch(name, id string) bool {
	for _, user := range UserMap {
		if user.StudentID == id && user.UserName == name {
			return true
		}
	}
	return false

}
func Login(res http.ResponseWriter, req *http.Request) {

	student_id, ok := alreadyLoggedIn(req)
	if student_id != "" && ok {

		// if already login in, switch to login_confirm page.
		http.Redirect(res, req, "/ExamPage/"+student_id, http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		username := sanitize.HTML(req.FormValue("UserName"))
		student_id := sanitize.HTML(req.FormValue("StudentID"))

		// check if user exist with username
		ok := CheckUserInfoMatch(username, student_id)
		// fmt.Println("match user name and student id", ok)
		if !ok {
			http.Error(res, "Username or student id does not match.", http.StatusUnauthorized)
			fmt.Println("not match!")
			return
		}

		// create session
		id, _ := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  COOKIE_NAME,
			Value: id.String(),
		}

		http.SetCookie(res, myCookie)
		// mu_session.Lock()
		mapSessions[myCookie.Value] = student_id

		fmt.Println("update login txt file and student id:", student_id)
		// check and update the logging time. If it is the first time logging,
		// record the login time and update the log file, otherwise, logging time
		// uses the earlier one.

		login_info, ok := LoginMap[student_id]
		fmt.Println("login info:", login_info)
		if !ok {
			fmt.Println("Student is the first time login.")
		}
		if login_info.LoggingTime == "" {
			login_info.LoggingTime = time.Now().Format("2006-01-02 15:04:05")
			SaveToTestLogFile()
		}
		// fmt.Println(*login_info)
		fmt.Println("redirect...")
		// mu_session.Unlock()
		http.Redirect(res, req, "/ExamPage/"+student_id, http.StatusSeeOther)
		return
	}

	// without passing any info to the login.gohtml, the page will show login failure information
	tpl.ExecuteTemplate(res, "login.html", nil)
}

// delete session cookie, then log out
func Submitted(res http.ResponseWriter, req *http.Request) {
	myCookie, err := req.Cookie(COOKIE_NAME)
	if err != nil {
		tpl.ExecuteTemplate(res, "submitted.gohtml", nil)
		return
	}
	// delete the session. the MaxAge is set to -1, which will override the
	// value of expires.
	// mu_session.Lock()
	delete(mapSessions, myCookie.Value)
	// mu_session.Unlock()

	// remove the cookie
	myCookie = &http.Cookie{
		Name:   COOKIE_NAME,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	tpl.ExecuteTemplate(res, "submitted.gohtml", nil)

}

// check if the session cookie has the userID. if it has, means already login in. Keep the current operation.
func alreadyLoggedIn(req *http.Request) (string, bool) {
	myCookie, err := req.Cookie(COOKIE_NAME)
	if err != nil {
		// that means the cookie does not exist
		return "", false
	}

	// mu_session.RLock()
	student_id, exist := mapSessions[myCookie.Value] //according to the uuid user session, get User object
	// mu_session.RUnlock()

	if !exist {
		// If the session id is not present in session map, that means the
		// user session has been terminated by server side.
		return "", false
	}

	fmt.Println(student_id)
	// if the user has been logged in, return user_id
	return student_id, true

}
