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

	var conect *sql.DB

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		//caso não haja o arquivo .env não esteja presente
		conect = conectarDBHeroku()
	} else {
		conect = conectarDBLocal()
	}

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	panic("PORT not defined")
	// }

	//com conect estamos instanciando a func conectarDB que sera passada como argumento do handler Usuario(*sql.DB)
	//conect := conectarDB()

	//aqui chamamos a func seed() para migrar os dados do []UsuariosDB para Banco de Dados novo.
	//depois que os dados foram migrados, podem deixar de chamar a função seed()
	//seed(conect)

	//temos 2 handlers: Home e Usuario
	http.HandleFunc("/", Home(conect))
	http.HandleFunc("/usuario/", Usuario(conect))

	s := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       10 * time.Minute,
		WriteTimeout:      2 * time.Minute,
		MaxHeaderBytes:    1 << 20,
	}
	log.Fatal(s.ListenAndServe())

}

//Home é uma função que vai usar o Template index.html e injeta informações de usuarios em uma tabela
func Home(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		logging(r)

		rows, err := db.Query("SELECT id, nome, sobrenome, email, perfil, mandato, foto, naturalidade FROM usuarios;")
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		linhas := make([]usuarios.Usuarios, 0)
		for rows.Next() {
			linha := usuarios.Usuarios{}
			err := rows.Scan(&linha.ID, &linha.Nome, &linha.Sobrenome, &linha.Email, &linha.Perfil, &linha.Mandato, &linha.Foto, &linha.Naturalidade)
			if err != nil {
				panic(err)
			}
			linhas = append(linhas, linha)
		}

		var tpl *template.Template
		tpl = template.Must(template.ParseGlob("./templates/*.html"))
		err = tpl.ExecuteTemplate(w, "index.html", linhas)
		if err != nil {
			panic(err)
		}
	}

}

//Usuario vai retornar um usuario que tenha o mesmo ID informado no http request
func Usuario(db *sql.DB) http.HandlerFunc {

	//retornamos um Handler será uma função anônima
	return func(w http.ResponseWriter, r *http.Request) {

		logging(r)

		//"params" é o URL de request. Nesse caso, /usuario/{id}
		params := r.URL.Path
		//"id" é o params sem /usuario/. Ficamos apenas com o numero que nos interessa: {id}
		id := strings.TrimPrefix(params, "/usuario/")
		//convertemos o tipo id de string para int e chamamos de "idint"
		idint, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("invalid param format")
		}

		//query armazena os dados do usuario que tenha ID igual ao numero informado o http request (idint)
		query := `SELECT id, nome, sobrenome, email, perfil, mandato, foto, naturalidade FROM usuarios WHERE id=$1;`

		//row terá o resultado da sql query
		row := db.QueryRow(query, idint)

		//criamos uma variável do tipo usuarios.Usuarios para receber as informações do banco de dados
		var usuario usuarios.Usuarios

		//copiamos o as informações de "row" para "usuario"
		err = row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Perfil, &usuario.Mandato, &usuario.Foto, &usuario.Naturalidade)
		if err != nil {
			fmt.Println(err)
		}

		usuarioSlice := make([]usuarios.Usuarios, 0)
		usuarioSlice = append(usuarioSlice, usuario)

		//Criamos um template tpl
		tpl := template.Must(template.ParseGlob("./templates/*.html"))
		//executamos o template com os dados presentes em "usuario" e enviamos o "response w"

		if len(usuarioSlice) == 0 {
			err = tpl.ExecuteTemplate(w, "usuarioNil.html", usuarioSlice)
			if err != nil {
				panic(err)
			}
		} else {
			err = tpl.ExecuteTemplate(w, "detalhesUsuario.html", usuarioSlice)
			if err != nil {
				panic(err)
			}
		}

	}
}

func logging(r *http.Request) {
	metodo := r.Method
	params := r.URL.Path
	host := r.Host
	fmt.Printf("Request: %s %s%s - %s\n", metodo, host, params, time.Now().Format(time.RFC3339))
}

//conectarDB vai fazer a interface entre o servidor e banco de dados usando as informações de acesso armazenadas no env do Heroku
func conectarDBHeroku() *sql.DB {
	DBinfo := fmt.Sprint(os.Getenv("DATABASE_URL"))
	db, err := sql.Open("pgx", DBinfo)
	if err != nil {
		panic(err)
	}

	return db
}

// conectarDBLocal vai fazer a interface entre o servidor e banco de dados usando as informações de acesso armazenadas no .env
func conectarDBLocal() *sql.DB {
	DBinfo := fmt.Sprintf("user=%s password=%s host=%s port=%v dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := sql.Open("pgx", DBinfo)
	if err != nil {
		panic(err)
	}

	return db
}

func seed(db *sql.DB) {
	//apaga tabela USUARIOS anterior, caso ela exista, e cria uma nova tabela com os campos abaixo
	query1 := `
	DROP TABLE IF EXISTS usuarios;
	CREATE TABLE usuarios (
		id SERIAL PRIMARY KEY,
		criado_em TIMESTAMP DEFAULT Now() NOT NULL,
		nome VARCHAR(50),
		sobrenome VARCHAR(50),
		email VARCHAR(100) UNIQUE,
		perfil VARCHAR(50),
		mandato VARCHAR(50),
		foto VARCHAR(350),
		naturalidade VARCHAR(100)
	);
	`
	_, err := db.Exec(query1)
	if err != nil {
		fmt.Println(err)
	}

	us := usuarios.UsuariosSlice

	for x := range us {
		usuario := us[x]
		query2 := `
		INSERT INTO usuarios(nome, sobrenome, email, perfil, mandato, foto, naturalidade)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
		_, err = db.Exec(query2, usuario.Nome, usuario.Sobrenome, usuario.Email, usuario.Perfil, usuario.Mandato, usuario.Foto, usuario.Naturalidade)
		if err != nil {
			panic(err)
		}
	}
}
