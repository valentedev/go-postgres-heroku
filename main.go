package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	http.HandleFunc("/usuario/", Usuario)

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

//Usuario é um handler que chama detalhesUsuario.html e mostra detalhes de cada usuario
func Usuario(w http.ResponseWriter, r *http.Request) {

	//"params" é o URL de request. Nesse caso, /usuario/{id}
	params := r.URL.Path
	//"id" é o params sem /usuario/. Ficamos apenas com o numero que nos interessa: {id}
	id := strings.TrimPrefix(params, "/usuario/")
	//convertemos o id de int para string e chamamos de "idint"
	idint, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("invalid param format")
	}
	//"us" é uma instância dos valores de Usuarios
	us := usuarios.UsuariosSlice
	//"resultado" é a função "encontrar" com os parâmetros us e indint
	resultado := encontrar(us, idint)
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

//função para ser usada no handler Usuario que usa o param do URL de request para encontrar valor dentro de um objeto da slice []usuarios.Usuarios
func encontrar(a []usuarios.Usuarios, b int) []usuarios.Usuarios {
	//criamos uma slice vazia chamada usuario (não confundir com Usuarios, no plural) que receberá o objeto encontrado
	var usuario = []usuarios.Usuarios{}
	for _, n := range a {
		if n.ID == b {
			//quando encontrado, n será adicionado ao slice "usuario"
			usuario = append(usuario, n)
		}
	}
	return usuario
}
