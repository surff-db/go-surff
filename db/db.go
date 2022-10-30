package surff

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const DB_FILE_NAME = "db.surff"

func DbCreate() {
	err := ioutil.WriteFile(DB_FILE_NAME, []byte(""), 0777)
	if err != nil {
		fmt.Println(err)
	}
}

func DbExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DbCheck() {
	if !DbExists(DB_FILE_NAME) {
		DbCreate()
	}
}

func Serialize(key string, value string) string {
	return "@" + key + "=" + value + ";"
}

func KeyExists(dbStr string, key string) string {
	r, _ := regexp.Compile(`(?m)@` + key + `=(.*?);`)
	value := r.FindString(dbStr)

	if value != "" {
		return value
	}

	return ""
}

func Get(key string) string {
	DbCheck()

	db, err := ioutil.ReadFile(DB_FILE_NAME)
	if err != nil {
		fmt.Println(err)
	}

	dbStr := string(db)

	r, _ := regexp.Compile(`(?m)@` + key + `=(.*?);`)
	value := r.FindStringSubmatch(dbStr)

	if value[1] != "" {
		return value[1]
	}

	return ""
}

func Set(key string, value string) {
	DbCheck()

	db, err := ioutil.ReadFile(DB_FILE_NAME)
	if err != nil {
		fmt.Println(err)
	}

	dbStr := string(db)

	replaceValue := KeyExists(dbStr, key)
	if replaceValue != "" {
		updateData := strings.Replace(dbStr, replaceValue, Serialize(key, value), -1)
		updateDataByte := []byte(updateData)

		err := ioutil.WriteFile(DB_FILE_NAME, updateDataByte, 0777)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		dbNew, err := os.OpenFile(DB_FILE_NAME, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer dbNew.Close()

		if _, err = dbNew.WriteString(Serialize(key, value)); err != nil {
			panic(err)
		}
	}
}

func Del(key string) {
	DbCheck()

	db, err := ioutil.ReadFile(DB_FILE_NAME)
	if err != nil {
		fmt.Println(err)
	}

	dbStr := string(db)

	replaceValue := KeyExists(dbStr, key)
	if replaceValue != "" {
		updateData := strings.Replace(dbStr, replaceValue, "", -1)
		updateDataByte := []byte(updateData)

		err := ioutil.WriteFile(DB_FILE_NAME, updateDataByte, 0777)
		if err != nil {
			fmt.Println(err)
		}
	}
}
