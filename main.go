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

	port := os.Getenv("PORT")

	//com conect estamos instanciando a func conectarDB que sera passada como argumento do handler Usuario(*sql.DB)
	//conect := conectarDB()

	//aqui chamamos a func seed() para migrar os dados do []UsuariosDB para Banco de Dados novo.
	//depois que os dados foram migrados, podem deixar de chamar a função seed()
	//seed(conect)

	// Usando http.Server
	// s := &http.Server{
	// 	Addr:              ":" + port,
	// 	ReadHeaderTimeout: 20 * time.Second,
	// 	ReadTimeout:       10 * time.Minute,
	// 	WriteTimeout:      2 * time.Minute,
	// 	MaxHeaderBytes:    1 << 20,
	// }
	// log.Fatal(s.ListenAndServe())

	//handlers funcs
	// http.HandleFunc("/", Home(conect))
	// http.HandleFunc("/usuario/", Usuario(conect))
	// http.HandleFunc("/usuario/criar/", CriarUsuario())
	// http.HandleFunc("/usuario/criado/", NovoUsuarioConfirma(conect))

	// Usuando http.ServerMux
	mux := http.NewServeMux()
	//handle do /static/
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", Home(conect))
	mux.HandleFunc("/usuario/", Usuario(conect))
	mux.HandleFunc("/usuario/criar/", CriarUsuario())
	mux.HandleFunc("/usuario/criado/", NovoUsuarioConfirma(conect))
	addr := ":" + port
	err = http.ListenAndServe(addr, mux)
	log.Fatal(err)

}

//Home é uma função que vai usar o Template index.html e injeta informações de usuarios em uma tabela
func Home(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		logging(r)

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		rows, err := db.Query("SELECT id, nome, sobrenome, email, perfil, mandato, foto, naturalidade FROM usuarios ORDER BY id DESC;")
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

		tpl = template.Must(template.ParseGlob("./templates/*"))
		//aqui passamos o nome do TEMPLATE e não do arquivo - nesse caso Index
		err = tpl.ExecuteTemplate(w, "Index", linhas)
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
		tpl := template.Must(template.ParseGlob("./templates/*"))
		//executamos o template com os dados presentes em "usuario" e enviamos o "response w"

		if len(usuarioSlice) == 0 {
			err = tpl.ExecuteTemplate(w, "usuarioNil.html", usuarioSlice)
			if err != nil {
				panic(err)
			}
		} else {
			err = tpl.ExecuteTemplate(w, "Detalhes", usuarioSlice)
			if err != nil {
				panic(err)
			}
		}

	}
}

//CriarUsuario gera um formulário para entrada de dados no DB
func CriarUsuario() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging(r)

		var tpl *template.Template
		tpl = template.Must(template.ParseGlob("./templates/*"))
		err := tpl.ExecuteTemplate(w, "Novo", nil)
		if err != nil {
			panic(err)
		}
	}
}

//NovoUsuarioConfirma faz o Parse da informação gerada em CriarUsuario() e inclui usuario no DB
func NovoUsuarioConfirma(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logging(r)

		//caso o método do request não seja POST, redireciona para o formulário de criação do usuario
		if r.Method != http.MethodPost {
			http.Error(w, "Método não autorizado", 405)
			//http.Redirect(w, r, "/usuario/criar/", 303)
			return
		}

		//instancia valores enviada pelo formulário
		nome := r.FormValue("nome")
		sobrenome := r.FormValue("sobrenome")
		email := r.FormValue("email")
		perfil := r.FormValue("perfil")
		mandato := r.FormValue("mandato")
		foto := r.FormValue("foto")
		naturalidade := r.FormValue("naturalidade")
		if nome == "" || sobrenome == "" || email == "" || perfil == "" || mandato == "" || foto == "" || naturalidade == "" {
			http.Redirect(w, r, "/usuario/criar/", http.StatusSeeOther)
		}

		query := `INSERT INTO usuarios (nome, sobrenome, email, perfil, mandato, foto, naturalidade) VALUES ($1,$2,$3,$4,$5,$6,$7);`

		_, err := db.Exec(query, nome, sobrenome, email, perfil, mandato, foto, naturalidade)
		if err != nil {
			panic(err)
		}

		//criamos uma variável do tipo usuarios.Usuarios para receber as informações do banco de dados
		var usuario usuarios.Usuarios

		//query armazena os dados do usuario que tenha ID igual ao numero informado o http request (idint)
		query = `SELECT id, nome, sobrenome, email, perfil, mandato, foto, naturalidade FROM usuarios WHERE email=$1;`

		//row terá o resultado da sql query
		row := db.QueryRow(query, email)

		//copiamos o as informações de "row" para "usuario"
		err = row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Perfil, &usuario.Mandato, &usuario.Foto, &usuario.Naturalidade)
		if err != nil {
			fmt.Println(err)
		}

		usuarioSlice := make([]usuarios.Usuarios, 0)
		usuarioSlice = append(usuarioSlice, usuario)

		var tpl *template.Template
		tpl = template.Must(template.ParseGlob("./templates/*"))
		err = tpl.ExecuteTemplate(w, "Confirma", usuarioSlice)
		if err != nil {
			panic(err)
		}
	}
}

//EditarUsuario é um handler para editar usuarios
func EditarUsuario(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Manda um log para o stdout
		logging(r)

		//Verificar seu método http é POST
		if r.Method != http.MethodPost {
			http.Error(w, "Método não autorizado", 405)
			//http.Redirect(w, r, "/usuario/criar/", 303)
			return
		}

		//Chama o template para edição do registro
		var tpl *template.Template
		tpl = template.Must(template.ParseGlob("./templates/*.html"))
		err := tpl.ExecuteTemplate(w, "formEditarUsuario.html", nil)
		if err != nil {
			panic(err)
		}

		// //Aloca o param ao var "id" e converte para integer
		// params := r.URL.Path
		// id := strings.TrimPrefix(params, "/usuario/editar")
		// idint, err := strconv.Atoi(id)
		// if err != nil {
		// 	fmt.Println("invalid param format")
		// }

		// query := `SELECT id, nome, sobrenome, email, perfil, mandato, foto, naturalidade FROM usuarios WHERE id=$1;`
		// row := db.QueryRow(query, idint)
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
