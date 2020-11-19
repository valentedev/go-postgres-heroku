package main

import (
	"fmt"
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
	fmt.Println(usuariosSlice)
	fmt.Printf("Numero de Usuários: %v", len(usuariosSlice))
	log.Fatal(s.ListenAndServe())

}

//Usuarios é uma struct que servirá como base para o futuro banco de dados
type Usuarios struct {
	ID        int
	Nome      string
	Sobrenome string
	Email     string
	Perfil    string
}

//usuariosSlice vai reunir um ou mais Usuários que serão renderizados no Template index.html
var usuariosSlice = []Usuarios{
	{
		ID:        1,
		Nome:      "Rodrigo",
		Sobrenome: "Valente",
		Email:     "valentedev.rodrigo@gmail.com",
		Perfil:    "admin",
	},
	{
		ID:        2,
		Nome:      "Campos",
		Sobrenome: "Sales",
		Email:     "cs@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        3,
		Nome:      "Rodrigues",
		Sobrenome: "Alves",
		Email:     "ra@gmail.com",
		Perfil:    "no-admin",
	}, {
		ID:        4,
		Nome:      "Afonso",
		Sobrenome: "Pena",
		Email:     "ap@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        5,
		Nome:      "Nilo",
		Sobrenome: "Peçanha",
		Email:     "np@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        6,
		Nome:      "Hermes",
		Sobrenome: "da Fonseca",
		Email:     "hf@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        7,
		Nome:      "Venceslau",
		Sobrenome: "Brás",
		Email:     "vb@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        8,
		Nome:      "Delfim",
		Sobrenome: "Moreira",
		Email:     "dm@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        9,
		Nome:      "Epitácio",
		Sobrenome: "Pessoa",
		Email:     "ra@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        10,
		Nome:      "Artur",
		Sobrenome: "Bernardes",
		Email:     "ab@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        11,
		Nome:      "Washington",
		Sobrenome: "Luis",
		Email:     "wl@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        12,
		Nome:      "Julio",
		Sobrenome: "Prestes",
		Email:     "jp@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        13,
		Nome:      "Getulio",
		Sobrenome: "Vargas",
		Email:     "gv@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        14,
		Nome:      "Jose",
		Sobrenome: "Linhares",
		Email:     "jl@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        15,
		Nome:      "Eurico",
		Sobrenome: "Gaspar Dutra",
		Email:     "egd@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        16,
		Nome:      "Café",
		Sobrenome: "Filho",
		Email:     "cf@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        17,
		Nome:      "Carlos",
		Sobrenome: "Luz",
		Email:     "cl@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        18,
		Nome:      "Nereu",
		Sobrenome: "Ramos",
		Email:     "ra@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        19,
		Nome:      "Juscelino",
		Sobrenome: "Kubitscheck",
		Email:     "jk@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        20,
		Nome:      "Jânio",
		Sobrenome: "Quadros",
		Email:     "jq@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        21,
		Nome:      "Ranieri",
		Sobrenome: "Mazzilli",
		Email:     "rm@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        22,
		Nome:      "João",
		Sobrenome: "Goulart",
		Email:     "jg@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        23,
		Nome:      "Humberto",
		Sobrenome: "Castelo Branco",
		Email:     "hcb@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        24,
		Nome:      "Artur",
		Sobrenome: "da Costa e Silva",
		Email:     "acs@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        25,
		Nome:      "Pedro",
		Sobrenome: "Aleixo",
		Email:     "pa@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        26,
		Nome:      "Emilio",
		Sobrenome: "Garrastazu Médici",
		Email:     "egm@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        27,
		Nome:      "Ernesto",
		Sobrenome: "Geisel",
		Email:     "eg@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        28,
		Nome:      "João",
		Sobrenome: "Figueiredo",
		Email:     "jf@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        29,
		Nome:      "José",
		Sobrenome: "Sarney",
		Email:     "js@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        30,
		Nome:      "Fernando",
		Sobrenome: "Collor",
		Email:     "fc@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        31,
		Nome:      "Itamar",
		Sobrenome: "Franco",
		Email:     "if@gmail.com",
		Perfil:    "no-admin",
	},
	{
		ID:        32,
		Nome:      "Fernando",
		Sobrenome: "Henrique Cardoso",
		Email:     "fhc@gmail.com",
		Perfil:    "no-admin",
	},
}

//Home é uma função que vai usar o Template index.html e injetar o usuario1
func Home(w http.ResponseWriter, r *http.Request) {

	var tpl *template.Template
	tpl = template.Must(template.ParseGlob("./templates/*.html"))
	err := tpl.ExecuteTemplate(w, "index.html", usuariosSlice)
	if err != nil {
		panic(err)
	}
}
