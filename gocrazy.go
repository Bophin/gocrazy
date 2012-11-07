package main

import (
//	"fmt"
	"net/http"
	"html/template"
)

var templ = template.Must(template.ParseFiles("webbstuff/GoCrazyIndex.html"))

func main() {
	http.HandleFunc("/", indexHandler)
	/*http.dir VERY UNSECURE*/
	http.Handle("/webbstuff/", http.StripPrefix("/webbstuff/", http.FileServer(http.Dir("webbstuff"))))
	http.ListenAndServe(":80", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("webbstuff/GoCrazyIndex.html")
	t.Execute(w, nil)
}
