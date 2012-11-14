package main

import (
	"net/http"
)



func main() {
	http.HandleFunc("/", rootHandler)
	//Make sure nothing important lies in "main_web"
	http.Handle("/main_web/", http.StripPrefix("/main_web/", http.FileServer(http.Dir("main_web"))))
	http.ListenAndServe(":80", nil)
}

