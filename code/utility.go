package main

import (
	"crypto/rand"
	"html/template"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

func getProjectRootPath() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	return filepath.Dir(wd)
}

// ParseTemplates function could find the subdirectories under templates directory and
// parse them conveniently.
func ParseTemplates() *template.Template {
	path := getProjectRootPath()

	templ := template.New("")
	err := filepath.Walk(path+"/templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".gohtml") || strings.Contains(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})

	if err != nil {
		panic(err)
	}

	return templ
}

// CreateRandom will create a random number
func CreateRandom(maxValue int64) string {
	// create a random user_id
	nBig, err := rand.Int(rand.Reader, big.NewInt(maxValue))
	if err != nil {
		log.Panic("could not generate user id.")
	}

	return nBig.String()
}

// IsFileExist checks if a certain file exists or not.
func isFileExist(relativeFileName string) bool {

	_, err := os.Stat(relativeFileName)
	if err != nil {
		if os.IsNotExist(err) {

			log.Fatal("File not Found !!")
			return false
		}

	}
	return true
}

// check the current user is admin or not
func IsAdmin(user_name string) bool {

	// if AdminInfo.UserName == user_name {
	// 	return true
	// }
	return false
}
