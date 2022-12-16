package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func ReadFromLoginFile() {

	if !isFileExist(getProjectRootPath() + "/data/log_info.txt") {
		// log_info could be non-existing, because no student has login in.
		fmt.Println("log info file does not exist.")
		return
	}

	f, err := os.Open(getProjectRootPath() + "/data/log_info.txt")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal("Fail to open log_info.txt", err)
	}
	defer f.Close()
	fmt.Println("Successfully Opened log_info.txt")

	// read our opened file as a byte array.

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		s := strings.TrimSuffix(scanner.Text(), "\n")
		login_info := strings.Split(s, ",")

		fmt.Println(login_info)
		LoginMap[login_info[0]] = &LoginInfo{StudentID: login_info[0], LoggingTime: login_info[1],
			SubmittingTime: login_info[2], SubmittedFileName: login_info[3]}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(LoginMap)

}

func ReadFromUserFile() {

	if !isFileExist(getProjectRootPath() + "/data/student_info.txt") {
		log.Fatal("student info file does not exist.")
		return
	}

	f, err := os.Open(getProjectRootPath() + "/data/student_info.txt")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal("Fail to open student_info.txt", err)
	}
	defer f.Close()
	fmt.Println("Successfully Opened student_info.txt")

	// read our opened file as a byte array.

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		s := strings.TrimSuffix(scanner.Text(), "\n")
		user_info := strings.Split(s, ",")
		UserMap[user_info[1]] = &UserInfo{user_info[0], user_info[1]}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(UserMap)

}
