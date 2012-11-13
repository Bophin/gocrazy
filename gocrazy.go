package main

import (
	"fmt"
	"net/http"
	"html/template"
)

var templ = template.Must(template.ParseFiles("webbstuff/GoCrazyIndex.html"))
var login = template.Must(template.ParseFiles("webbstuff/login.html"))

func main() {
	http.HandleFunc("/", indexHandler)
	/*http.dir VERY UNSECURE*/
	http.Handle("/webbstuff/", http.StripPrefix("/webbstuff/", http.FileServer(http.Dir("webbstuff"))))
	http.HandleFunc("/login/", loginHandler);
	http.ListenAndServe(":80", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("webbstuff/GoCrazyIndex.html")
	t.Execute(w, nil)
}

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
