package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
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

	conect := conectarDB()

	http.HandleFunc("/", Home)
	http.HandleFunc("/usuario/", Usuario(conect))

	s := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       10 * time.Minute,
		WriteTimeout:      2 * time.Minute,
		MaxHeaderBytes:    1 << 20,
	}
	log.Fatal(s.ListenAndServe())

}

//Home é uma função que vai usar o Template index.html e injeta informações de usuarios em uma tabela
func Home(w http.ResponseWriter, r *http.Request) {

	var tpl *template.Template
	us := usuarios.UsuariosSlice
	tpl = template.Must(template.ParseGlob("./templates/*.html"))
	err := tpl.ExecuteTemplate(w, "index.html", us)
	if err != nil {
		panic(err)
	}
}

func Usuario(db *sql.DB) http.HandlerFunc {

	//ess é ema função anônima
	return func(w http.ResponseWriter, r *http.Request) {

		//"params" é o URL de request. Nesse caso, /usuario/{id}
		params := r.URL.Path
		//"id" é o params sem /usuario/. Ficamos apenas com o numero que nos interessa: {id}
		id := strings.TrimPrefix(params, "/usuario/")
		//convertemos o tipo id de string para int e chamamos de "idint"
		idint, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("invalid param format")
		}

		usdb := usuarios.UsuariosDB

		resultado := 

		//Criamos um template tpl
		tpl := template.Must(template.ParseGlob("./templates/*.html"))
		if len(resultado) == 0 {
			//se "resultado" retornar vazia executamos o template usuarioNil.html
			err = tpl.ExecuteTemplate(w, "usuarioNil.html", resultado)
			if err != nil {
				panic(err)
			}
		} else {
			//executamos o template detalhesUsuario.html
			err = tpl.ExecuteTemplate(w, "detalhesUsuario.html", resultado)
			if err != nil {
				panic(err)
			}
		}

	}
}



//conectarDB vai fazer a interface entre o servidor e banco de dados usando as informações de acesso armazenadas no .env
func conectarDB() *sql.DB {
	DBinfo := fmt.Sprintf("user=%s password=%s host=%s port=%v dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := sql.Open("pgx", DBinfo)
	if err != nil {
		panic(err)
	}

	return db
}
