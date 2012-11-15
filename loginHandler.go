package main

import (
	"fmt"
	"net/http"
	"html/template"
	"io/ioutil"
	"strings"
)

var logPath = "main_web/login.html"

func loginHandler(w http.ResponseWriter, r *http.Request) {

	t_login, err := template.ParseFiles(logPath)
	if err != nil {
		panic(err)
	}

	method := r.Method
	if method == "GET" {
		t_login.Execute(w, nil)
	} else if method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			t_login.Execute(w, nil)
		} else {
			if userCheck(username, password) {
				fmt.Fprintf(w, "You are great!")
			} else {
				fmt.Fprintf(w, "You are useless!")
			}
		}

	}
}

func userCheck(username, password string) bool {

	file, err := ioutil.ReadFile("content/login/pwd.txt")
	if err != nil {
		panic(err)
	}

	//Splits per row
	f_str := string(file)
	users := strings.Split(f_str, "\n")

	for _, val := range users {
		//Splits rows by :
		row := strings.Split(val, ":")
		if row[0] == username && row[1] == password {
			return true
		}
	}
	return false
}
