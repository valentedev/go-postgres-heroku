package main

import (
	"html/template"
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

//Usuarios é uma struct que servirá como base para o futuro banco de dados
type Usuarios struct {
	ID        int
	Nome      string
	Sobrenone string
	Email     string
	Perfil    string
}

//usuario1 uma instância de Usuarios
var usuario1 = Usuarios{
	ID:        1,
	Nome:      "Rodrigo",
	Sobrenone: "Valente",
	Email:     "valentedev.rodrigo@gmail.com",
	Perfil:    "admin",
}

//Home é uma função que vai usar o Template index.html e injetar o usuario1
func Home(w http.ResponseWriter, r *http.Request) {

	var tpl *template.Template
	tpl = template.Must(template.ParseGlob("./templates/*.html"))
	err := tpl.ExecuteTemplate(w, "index.html", usuario1)
	if err != nil {
		panic(err)
	}
}
