package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/valentedev/go-postgres-heroku/usuarios"
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

	http.HandleFunc("/", Home)

	s := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       10 * time.Minute,
		WriteTimeout:      2 * time.Minute,
		MaxHeaderBytes:    1 << 20,
	}
	log.Fatal(s.ListenAndServe())

}

//Home é uma função que vai usar o Template index.html e injetar o usuario1
func Home(w http.ResponseWriter, r *http.Request) {

	var tpl *template.Template
	us := usuarios.UsuariosSlice
	tpl = template.Must(template.ParseGlob("./templates/*.html"))
	err := tpl.ExecuteTemplate(w, "index.html", us)
	if err != nil {
		panic(err)
	}
}
