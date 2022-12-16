package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/kennygrant/sanitize"
	"golang.org/x/crypto/bcrypt"
)

func CheckUserInfoMatch(name, id string) bool {
	for _, user := range UserMap {
		if user.StudentID == id && user.UserName == name {
			return true
		}
	}
	return false

}

func checkAdminInfoMatch(input_user_name, input_password string) bool {

	if strings.EqualFold(Admin.UserName, input_user_name) {
		err := bcrypt.CompareHashAndPassword([]byte(Admin.Password), []byte(input_password))
		if err != nil {
			return false
		}

	}
	return true

}

func AdminLogin(res http.ResponseWriter, req *http.Request) {

	user_id, ok := alreadyLoggedIn(req)
	if user_id != "" && ok {

		// if already login in, switch to login_confirm page.
		// if already login in, switch to login_confirm page.
		http.Redirect(res, req, "/admin/setting", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		username := sanitize.HTML(req.FormValue("UserName"))
		password := sanitize.HTML(req.FormValue("Password"))

		// check if user exist with username
		ok := checkAdminInfoMatch(username, password)
		if !ok {
			http.Error(res, "Username or password does not match.", http.StatusUnauthorized)
			return
		}

		// create session
		id, _ := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  COOKIE_NAME,
			Value: id.String(),
		}

		http.SetCookie(res, myCookie)
		mapSessions[myCookie.Value] = user_id
		http.Redirect(res, req, "/admin/setting", http.StatusSeeOther)
		return
	}

	// without passing any info to the login.gohtml, the page will show login failure information
	tpl.ExecuteTemplate(res, "admin_login.html", nil)
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
