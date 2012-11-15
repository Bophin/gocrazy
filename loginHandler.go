package main

import (
	"fmt"
	"net/http"
	"html/template"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	t_login, err := template.ParseFiles(files + "login.html")
	if err != nil {
		panic(err)
	}

	method := r.Method
	if method == "GET" {
		t_login.Execute(w, nil)
	} else if method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "olov" && password == "jesper" {
			fmt.Fprintf(w, "You are great!")
		} else {
			fmt.Fprintf(w, "You are useless!")
		}
	}
}
