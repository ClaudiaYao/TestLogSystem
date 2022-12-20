package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/kennygrant/sanitize"
)

func ExamPage(res http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)

	if req.Method == http.MethodGet {

		// not logged in.
		login_student_id, err := alreadyLoggedIn(req)
		if !err {
			http.Redirect(res, req, "/login", http.StatusSeeOther)
			return
		}

		// check if the studentID passed as query parameter exists or not.
		student_id := sanitize.HTML(params["StudentID"])
		if login_student_id != student_id {
			fmt.Fprintln(res, "You could not see other's result.")
			return
		}
		// if it is the login user who wants to enter the Exam Page
		_, ok := UserMap[student_id]
		if !ok {
			http.Error(res, "Username or student id does not match.", http.StatusUnauthorized)
			fmt.Println("not match!")
			return

		} else {

			// match user info, need to check if it is already submitted. Could not
			// enter the Exam Page after submitting.
			user := UserMap[student_id]
			user_login := LoginMap[student_id]
			fmt.Println(user_login)

			// already submitted
			if user_login.SubmittedFileName != "" {
				fmt.Println("has already submitted. Could not re-submit again.")
				http.Redirect(res, req, "/submitted", http.StatusSeeOther)
				return
			}

			// not submitted, then check the starting time
			fmt.Println(user_login.LoggingTime)
			start_time, error := time.Parse("2006-01-02 15:04:05", user_login.LoggingTime)

			if error != nil {
				fmt.Println(error)
				return
			}

			// if already logged before, add two hours as the final ending time
			// when the user first logins, the starting time has already been logged,
			// so here we do not need to consider the missing starting time condition.
			end_time := start_time.Add(time.Hour * 2)
			user_info := map[string]interface{}{"UserName": user.UserName, "StudentID": user.StudentID,
				"StartTime": user_login.LoggingTime, "EndTime": end_time}
			tpl.ExecuteTemplate(res, "examPage.html", user_info)
		}

	} else if req.Method == http.MethodPost {
		fmt.Println("submit the files.")
		student_id := sanitize.HTML(params["StudentID"])
		fmt.Println("posting: get student id:", student_id)
		fmt.Println(req.FormValue("Submit"))
		// when user submitted files, call the function uploadFiles, also
		// pass the student id, so that the zipped folder will use student id info.
		if req.FormValue("Submit") != "" {
			// the returned short file_name will be used to update log_info.txt
			file_name, ok := uploadFile(res, req, student_id)
			fmt.Println(file_name, ok)
			if !ok {
				fmt.Println("fail to submit your test paper. Resubmit again.")
			} else {
				userlog := LoginMap[student_id]
				userlog.SubmittedFileName = file_name
				userlog.SubmittingTime = time.Now().Format("2006-01-02 15:04:05")
				fmt.Println(UserMap[student_id])
				SaveToLogFile()

				http.Redirect(res, req, "/submitted", http.StatusSeeOther)
			}
		}

	}
}

func uploadFile(w http.ResponseWriter, r *http.Request, studentID string) (string, bool) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("UploadFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return "", false
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	user := UserMap[studentID]
	file_name := studentID + "_" + user.UserName + "_" + handler.Filename
	// create a file under the specific directory
	save_file, err := os.Create(getProjectRootPath() + "/Submission/" + file_name)
	// tempFile, err := ioutil.TempFile("Submission", studentID+"_"+user.UserName+"_*.zip")
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	defer save_file.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	// write this byte array to our temporary file
	save_file.Write(fileBytes)
	// return that we have successfully uploaded our file!

	return getShortFileName(save_file.Name()), true
	// fmt.Fprintf(w, "Successfully Uploaded File\n")
}


