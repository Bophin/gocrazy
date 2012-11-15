package main

import (
	"net/http"
	"flag"
)

var port = flag.String("p", "80", "Sets server port")

func main() {
	flag.Parse()

	http.HandleFunc("/", rootHandler)
	//Make sure nothing important lies in "main_web"
	http.Handle("/main_web/", http.StripPrefix("/main_web/", http.FileServer(http.Dir("main_web"))))
	http.HandleFunc("/login/", loginHandler);
	http.ListenAndServe(":" + *port, nil)
}

