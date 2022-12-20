package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 10 // 10MB

func AdminSetting(res http.ResponseWriter, req *http.Request) {

	user_name, ok := alreadyLoggedIn(req)
	if !ok {

		// if not login in, switch to login page.
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		AdminUpload(res, req)
	} else if req.Method == http.MethodGet {
		tpl.ExecuteTemplate(res, "settings.html", user_name)
	}

}

// Admin could upload multiple files
func AdminUpload(res http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req.Body = http.MaxBytesReader(res, req.Body, MAX_UPLOAD_SIZE)
	if err := req.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(res, "The uploaded file is too big. Please choose an file that's less than 10 MB in size", http.StatusBadRequest)
		return
	}

	formdata := req.MultipartForm // ok, no problem so far, read the Form data
	//get the *fileheaders
	files := formdata.File["UploadFile"] // grab the filenames

	for i := range files { // loop through the files one by one
		fileHeader := files[i]
		file, err := fileHeader.Open()

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// since req.FormFile only returns the first file for the provided form key,
		// so we need to use the other method.
		// file, fileHeader, err := req.FormFile("UploadFile")
		// if err != nil {
		// 	http.Error(res, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// defer file.Close()

		// Create the uploads folder if it doesn't
		// already exist
		upload_folder := getProjectRootPath() + "/admin_uploads"
		err = os.MkdirAll(upload_folder, os.ModePerm)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a new file in the uploads directory
		dst, err := os.Create(upload_folder + "/" + getShortFileName(fileHeader.Filename))
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		// Copy the uploaded file to the filesystem
		// at the specified destination
		_, err = io.Copy(dst, file)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(res, "Upload successful")
	}
}
