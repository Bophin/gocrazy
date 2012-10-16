package main

import (
	"fmt"
	"net/http"
	"html/template"
)

var templ = template.Must(template.ParseFiles("webbstuff/index.html"))

func main() {
	fmt.Println("Hej!")
	fmt.Println("Hej världen! Jag är med <3")
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templ.Execute(w, nil)
}
