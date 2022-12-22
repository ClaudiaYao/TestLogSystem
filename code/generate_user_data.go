package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	// "golang.org/x/crypto/bcrypt"
)

var TestUserMap = make(map[string]*UserInfo)
var TestLogMap = make(map[string]*LoginInfo)
var TestAdmin AdminInfo

// this part is just for creating mock data. The map stores structs, not the pointer to structs
func generate_users() {

	first_names := []string{"James", "Robert", "John", "Michael", "David", "William", "Richard",
		"Joseph", "Thomas", "Charles", "Christopher", "Daniel", "Matthew", "Anthony", "Mark",
		"Donald", "Steven", "Paul", "Andrew", "Joshua", "Kenneth", "Kevin", "Brian", "George",
		"Timothy", "Patricia", "Jennifer", "Linda", "Elizabeth", "Barbara", "Susan",
		"Jessica", "Sarah", "Karen", "Lisa", "Nancy", "Betty", "Margaret", "Sandra",
		"Ashley", "Kimberly", "Emily", "Donna", "Michelle", "Carol", "Amanda", "Dorothy",
		"Melissa", "Deborah", "Stephanie", "Rebecca"}
	len_first := len(first_names)

	last_names := []string{"Tan", "Lim", "Lee", "Ng", "Ong", "Wong", "Goh",
		"Chua", "Chan", "Koh", "Teo", "Ang", "Yeo", "Tay", "Ho", "Low", "Toh", "Sim",
		"Chong", "Chia", "Seah"}
	len_last := len(last_names)

	for count := 0; count < 100; count++ {
		student_id := "ST" + strconv.Itoa(70000000+rand.Intn(10000000))
		first_name := first_names[rand.Intn(len_first)]
		last_name := last_names[rand.Intn(len_last)]
		user_name := first_name + last_name

		TestUserMap[student_id] = &UserInfo{user_name, student_id}

	}
	SaveToUserFile()
}

func generate_logging_data() {

	for k := range TestUserMap {
		if rand.Intn(5) > 2 {

			random_time := time.Now().Add(time.Minute * time.Duration(rand.Intn(30)))
			login_time := random_time.Format("2006-01-02 15:04:05")

			random_time_2 := random_time.Add(time.Minute * time.Duration(rand.Intn(60)))
			submitted_time := random_time_2.Format("2006-01-02 15:04:05")

			TestLogMap[k] = &LoginInfo{k, login_time, submitted_time, ""}

		} else if rand.Intn(3) > 2 {
			TestLogMap[k] = &LoginInfo{StudentID: k}
		} else {
			random_time := time.Now().Add(-time.Minute * time.Duration(rand.Intn(30)))
			login_time := random_time.Format("2006-01-02 15:04:05")
			TestLogMap[k] = &LoginInfo{StudentID: k, LoggingTime: login_time}
		}

	}

	SaveToTestLogFile()

}

// This function will be called each time when a user session finishes
// just to back up the record in time. Since we do not use database
// in this case, we need to keep updating and maintaining the json files.
func SaveToTestLogFile() {

	fmt.Println("write to log_info.txt")
	f, err := os.Create(getProjectRootPath() + "/data/log_info.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for _, loginlog := range TestLogMap {
		// user := TestUserMap[k]
		_, err2 := f.WriteString(loginlog.StudentID + "," + loginlog.LoggingTime + "," +
			loginlog.SubmittingTime + "," + loginlog.SubmittedFileName + "\n")

		if err2 != nil {
			log.Println("fail to save to log_info.txt", err)
			log.Fatal(err2)
		}
	}

}

func SaveToUserFile() {

	fmt.Println("write to user_info.txt")
	f, err := os.Create(getProjectRootPath() + "/data/student_info.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for _, user := range TestUserMap {

		_, err2 := f.WriteString(user.UserName + "," + user.StudentID + ",\n")

		if err2 != nil {
			log.Println("fail to save to appointment.json: ", err)
			log.Fatal(err2)
		}
	}

}

func generate_admin_data() {

	pwd, err := bcrypt.GenerateFromPassword([]byte("Admin"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("fail to generate admin password")
		log.Panic("fail to generate admin password")
	}

	TestAdmin = AdminInfo{UserName: "Claudia", Password: pwd}
	// fmt.Println(bcrypt.GenerateFromPassword([]byte(input_password), bcrypt.DefaultCost))

	SaveToAdminJSON()

}

func SaveToAdminJSON() {

	json_file, _ := json.MarshalIndent(TestAdmin, "", " ")

	err := os.WriteFile(getProjectRootPath()+"/data/admin_info.json", json_file, 0644)
	if err != nil {
		log.Println("fail to save to admin_info.json error: ", err)
	}

}

// This function will be called each time when a user session finishes
// just to back up the record in time. Since we do not use database
// in this case, we need to keep updating and maintaining the json files.
func SaveToTXTFiles() {
	generate_users()
	generate_logging_data()

}
