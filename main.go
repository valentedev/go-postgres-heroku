package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Olá, %q", html.EscapeString(r.URL.Path))
	})

	s := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       10 * time.Minute,
		WriteTimeout:      2 * time.Minute,
		MaxHeaderBytes:    1 << 20,
	}
	log.Fatal(s.ListenAndServe())

}
