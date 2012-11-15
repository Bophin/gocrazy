package main

import (
	"fmt"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "GET" {
		login.Execute(w, nil)
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
