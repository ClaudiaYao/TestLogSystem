package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func ReadJsonAdmin() {

	if !isFileExist(getProjectRootPath() + "/data/admin_info.json") {
		return
	}
	jsonFileAdmin, err := os.Open(getProjectRootPath() + "/data/admin_info.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Panic("Fail to open admin_info.json", err)
	}

	defer jsonFileAdmin.Close()
	fmt.Println("Successfully Opened admin_info.json")
	// defer the closing of our jsonFile so that we can parse it later on

	byteValue, _ := ioutil.ReadAll(jsonFileAdmin)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &Admin)
	fmt.Println(Admin)

}
