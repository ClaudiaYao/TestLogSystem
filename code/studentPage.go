package main

import (
	"fmt"
	"net/http"
	"time"

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
