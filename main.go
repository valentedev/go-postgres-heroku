package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT not defined")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Olá, %q", html.EscapeString(r.URL.Path))
	})

	s := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       10 * time.Minute,
		WriteTimeout:      2 * time.Minute,
		MaxHeaderBytes:    1 << 20,
	}
	log.Fatal(s.ListenAndServe())

}
