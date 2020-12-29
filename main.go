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

	"github.com/dgrijalva/jwt-go"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/valentedev/go-postgres-heroku/usuarios"
	"golang.org/x/crypto/bcrypt"
)

//com conect estamos instanciando a func conectarDB que sera passada como argumento do handler(*sql.DB)
var conect *sql.DB

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		//caso não haja o arquivo .env não esteja presente
		conect = conectarDBHeroku()
	} else {
		conect = conectarDBLocal()
	}

	port := os.Getenv("PORT")

	// Usuando http.ServerMux
	mux := http.NewServeMux()
	//handle do /static/
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/login/", Login(conect))
	mux.HandleFunc("/login/redirect/", Logado(conect))
	mux.HandleFunc("/", TokenMiddleware(Home(conect)))
	mux.HandleFunc("/usuario/", TokenMiddleware(Usuario(conect)))
	mux.HandleFunc("/usuario/criar/", TokenMiddleware(CriarUsuario()))
	mux.HandleFunc("/usuario/criado/", TokenMiddleware(Criado(conect)))
	mux.HandleFunc("/usuario/editar/", TokenMiddleware(EditarUsuario(conect)))
	mux.HandleFunc("/usuario/editado/", TokenMiddleware(Editado(conect)))
	mux.HandleFunc("/usuario/deletar/", TokenMiddleware(Deletar(conect)))
	mux.HandleFunc("/usuario/deletado/", TokenMiddleware(Deletado(conect)))

	// // aqui chamamos a func seed() para migrar os dados do []UsuariosDB para Banco de Dados novo.
	// // depois que os dados foram migrados, podem deixar de chamar a função seed(db *sql.DB)
	//seed(conect)

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

		rows, err := db.Query("SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios ORDER BY id DESC;")
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		linhas := make([]usuarios.Usuarios, 0)
		for rows.Next() {
			linha := usuarios.Usuarios{}
			err := rows.Scan(&linha.ID, &linha.Nome, &linha.Sobrenome, &linha.Email, &linha.Senha, &linha.Admin, &linha.Ativo)
			if err != nil {
				panic(err)
			}
			linhas = append(linhas, linha)
		}

		type Dados struct {
			Linhas []usuarios.Usuarios
			//Usuario string
		}

		dados := Dados{
			Linhas: linhas,
			//Usuario: t,
		}

		var tpl *template.Template

		tpl = template.Must(template.ParseGlob("./templates/*"))

		err = tpl.ExecuteTemplate(w, "Index", dados)
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
		query := `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE id=$1;`

		//row terá o resultado da sql query
		row := db.QueryRow(query, idint)

		//criamos uma variável do tipo usuarios.Usuarios para receber as informações do banco de dados
		var usuario usuarios.Usuarios

		//copiamos o as informações de "row" para "usuario"
		err = row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			fmt.Println(err)
		}

		//Criamos um template tpl
		tpl := template.Must(template.ParseGlob("./templates/*"))
		//executamos o template com os dados presentes em "usuario" e enviamos o "response w"
		err = tpl.ExecuteTemplate(w, "Detalhes", usuario)
		if err != nil {
			panic(err)
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

//Criado faz o Parse da informação gerada em CriarUsuario() e inclui usuario no DB
func Criado(db *sql.DB) http.HandlerFunc {
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
		senha := r.FormValue("senha")
		admin := r.FormValue("admin")
		ativo := r.FormValue("ativo")
		if nome == "" || sobrenome == "" || email == "" {
			http.Redirect(w, r, "/usuario/criar/", http.StatusSeeOther)
		}

		senhaByte, err := SenhaHash(senha)
		if err != nil {
			fmt.Println(err)
		}

		senha = string(senhaByte)

		r.ParseForm()
		formAdmin := r.Form.Get("admin")
		if formAdmin == "" {
			admin = "false"
		}
		formAtivo := r.Form.Get("ativo")
		if formAtivo == "" {
			ativo = "false"
		}
		//fmt.Println(r.Form)

		query := `INSERT INTO usuarios (nome, sobrenome, email, senha, admin, ativo) VALUES ($1,$2,$3,$4,$5,$6);`

		_, err = db.Exec(query, nome, sobrenome, email, senha, admin, ativo)
		if err != nil {
			panic(err)
		}

		//criamos uma variável do tipo usuarios.Usuarios para receber as informações do banco de dados
		var usuario usuarios.Usuarios

		//query armazena os dados do usuario que tenha ID igual ao numero informado o http request (idint)
		query = `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE email=$1;`

		//row terá o resultado da sql query
		row := db.QueryRow(query, email)

		//copiamos o as informações de "row" para "usuario"
		err = row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			fmt.Println(err)
		}

		usuarioSlice := make([]usuarios.Usuarios, 0)
		usuarioSlice = append(usuarioSlice, usuario)

		var tpl *template.Template
		tpl = template.Must(template.ParseGlob("./templates/*"))
		err = tpl.ExecuteTemplate(w, "Criado", usuarioSlice)
		if err != nil {
			panic(err)
		}
	}
}

//EditarUsuario é um handler para editar usuarios
func EditarUsuario(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logging(r)

		//"params" é o URL de request. Nesse caso, /usuario/{id}
		params := r.URL.Path
		//"id" é o params sem /usuario/. Ficamos apenas com o numero que nos interessa: {id}
		id := strings.TrimPrefix(params, "/usuario/editar/")
		//convertemos o tipo id de string para int e chamamos de "idint"
		idint, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("invalid param format")
		}

		//query armazena os dados do usuario que tenha ID igual ao numero informado o http request (idint)
		query := `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE id=$1;`

		//row terá o resultado da sql query
		row := db.QueryRow(query, idint)

		//criamos uma variável do tipo usuarios.Usuarios para receber as informações do banco de dados
		var usuario usuarios.Usuarios

		//copiamos o as informações de "row" para "usuario"
		err = row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			fmt.Println(err)
		}

		var tpl *template.Template
		tpl = template.Must(template.ParseGlob("./templates/*"))
		err = tpl.ExecuteTemplate(w, "Editar", usuario)
		if err != nil {
			panic(err)
		}
	}

}

//Editado vai retornar um usuario que tenha o mesmo ID informado no http request
func Editado(db *sql.DB) http.HandlerFunc {

	//retornamos um Handler será uma função anônima
	return func(w http.ResponseWriter, r *http.Request) {

		logging(r)

		//caso o método do request não seja PUT, redireciona para o formulário de criação do usuario
		if r.Method != http.MethodPost {
			http.Error(w, "Método não autorizado", 405)
			return
		}

		//instancia valores enviada pelo formulário
		id := r.FormValue("id")
		idint, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}

		nome := r.FormValue("nome")
		sobrenome := r.FormValue("sobrenome")
		email := r.FormValue("email")
		senha := r.FormValue("senha")
		admin := r.FormValue("admin")
		ativo := r.FormValue("ativo")

		r.ParseForm()
		formAdmin := r.Form.Get("admin")
		if formAdmin == "" {
			admin = "false"
		}
		formAtivo := r.Form.Get("ativo")
		if formAtivo == "" {
			ativo = "false"
		}

		//fmt.Println(r.Form)

		//criamos uma variável do tipo usuarios.Usuarios para receber as informações do banco de dados
		var usuario usuarios.Usuarios

		//query armazena os dados do usuario que tenha ID igual ao numero informado o http request (idint)
		query := `UPDATE usuarios SET nome=$1, sobrenome=$2, email=$3, senha=$4, admin=$5, ativo=$6 WHERE id=$7;`

		// _, err = db.Exec(query, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Perfil, &usuario.Mandato, &usuario.Foto, &usuario.Naturalidade, idint)
		_, err = db.Exec(query, nome, sobrenome, email, senha, admin, ativo, idint)
		if err != nil {
			panic(err)
		}

		//query armazena os dados do usuario que tenha ID igual ao numero informado o http request (idint)
		query = `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE id=$1;`

		//row terá o resultado da sql query
		row := db.QueryRow(query, idint)

		//copiamos o as informações de "row" para "usuario"
		err = row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			fmt.Println(err)
		}

		//Criamos um template tpl
		tpl := template.Must(template.ParseGlob("./templates/*"))
		err = tpl.ExecuteTemplate(w, "Editado", usuario)
		if err != nil {
			panic(err)
		}
	}
}

// Deletar inicia o processo de remoção do usuário do banco de dados
func Deletar(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		logging(r)

		params := r.URL.Path
		id := strings.TrimPrefix(params, "/usuario/deletar/")
		idint, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("invalid param format")
		}

		query := `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE id=$1;`
		row := db.QueryRow(query, idint)
		var usuario usuarios.Usuarios
		err = row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			fmt.Println(err)
		}

		var tpl *template.Template
		tpl = template.Must(template.ParseGlob("./templates/*"))
		err = tpl.ExecuteTemplate(w, "Deletar", usuario)
		if err != nil {
			panic(err)
		}

	}

}

// Deletado confirma a remoção do usuário do banco de dados
func Deletado(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging(r)

		params := r.URL.Path
		id := strings.TrimPrefix(params, "/usuario/deletado/")
		idint, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("invalid param format")
		}

		query := `DELETE FROM usuarios WHERE id=$1;`
		_, err = db.Exec(query, idint)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/usuarios/", 307)

	}
}

// Login recebe email e senha e autentica acesso
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logging(r)

		var tpl *template.Template
		tpl = template.Must(template.ParseGlob("./templates/*"))
		err := tpl.ExecuteTemplate(w, "Login", nil)
		if err != nil {
			panic(err)
		}
	}
}

// Logado recebe email e senha fo Login e autentica acesso
func Logado(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logging(r)

		if r.Method != http.MethodPost {
			http.Error(w, "Método não autorizado", 405)
			return
		}

		email := r.FormValue("email")
		senha := r.FormValue("senha")
		if email == "" || senha == "" {
			http.Redirect(w, r, "/login/", http.StatusSeeOther)
		}

		// query baseada em email
		query := `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE email=$1`
		row := db.QueryRow(query, email)
		var usuario usuarios.Usuarios
		err := row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			fmt.Println(err)
		}

		token, err := Token(usuario)
		if err != nil {
			panic(err)
		}

		c := http.Cookie{
			Path:     "/",
			Name:     "session",
			Value:    token,
			HttpOnly: true,
			Expires:  time.Now().Add(time.Hour * 24),
			Secure:   false,
		}

		//nome := usuario.Nome

		if usuario.Admin == true {
			http.SetCookie(w, &c)
			//w.Header().Add("Authorization", token)
			//Email(nome)
			http.Redirect(w, r, "/", 307)
		} else {
			http.Error(w, "Acesso não autorizado", 401)
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
		nome VARCHAR(50) NOT NULL,
		sobrenome VARCHAR(50) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		senha VARCHAR(100),
		admin boolean DEFAULT false NOT NULL,
		ativo boolean DEFAULT false NOT NULL
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
		INSERT INTO usuarios(nome, sobrenome, email, senha, admin, ativo)
		VALUES ($1,$2,$3,$4,$5,$6)`
		_, err = db.Exec(query2, usuario.Nome, usuario.Sobrenome, usuario.Email, usuario.Senha, usuario.Admin, usuario.Ativo)
		if err != nil {
			panic(err)
		}
	}
}

// TOKEN #####################

type minhasClaims struct {
	jwt.StandardClaims
	Email string
	Nome  string
}

var assinatura = os.Getenv("TOKEN_SECRET")

// Token envia uma string para o client que será usada para autenticação.
func Token(u usuarios.Usuarios) (string, error) {

	var err error

	claims := minhasClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "go-postgres-heroku",
		},
		Email: u.Email,
		Nome:  u.Nome,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	tokenString, err := token.SignedString([]byte(assinatura))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

// TokenCheck verifica a validade do token
func TokenCheck(c *http.Cookie) (string, error) {

	tokenString := c.Value

	afterVerificationToken, err := jwt.ParseWithClaims(tokenString, &minhasClaims{}, func(beforeVeritificationToken *jwt.Token) (interface{}, error) {
		if beforeVeritificationToken.Method.Alg() != jwt.SigningMethodES256.Alg() {
			return nil, fmt.Errorf("Alguem tentou hackear o siging method")
		}
		return []byte(assinatura), nil
	})

	tokenOK := afterVerificationToken.Valid && err == nil

	mensagemAuth := "ninguém"
	claims := afterVerificationToken.Claims.(*minhasClaims)

	if tokenOK {
		mensagemAuth = claims.Nome
	}

	return mensagemAuth, nil
}

// TokenValid verifica se o token é válido e
func TokenValid(c *http.Cookie) (string, error) {
	tokenString := c.Value
	tokenVerificado, err := jwt.ParseWithClaims(tokenString, &minhasClaims{}, func(tokenNaoVerificado *jwt.Token) (interface{}, error) {
		return []byte(assinatura), nil
	})
	if !tokenVerificado.Valid {
		return "", err
	}
	return tokenString, nil
}

// TokenMiddleware é um wrapper que vai verificar se há um token válido em cada Handler.
func TokenMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("session")
		if err != nil {
			c = &http.Cookie{}
		}

		tokenString := c.Value

		tokenVerificado, err := jwt.ParseWithClaims(tokenString, &minhasClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(assinatura), nil
		})
		if err != nil || !tokenVerificado.Valid {
			fmt.Println("Token inválido ou inexistente")
			http.Redirect(w, r, "/login", 307)
		}

		next.ServeHTTP(w, r)
	})

}

// SenhaHash usa bcrypt para encriptar a senha
func SenhaHash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

// SenhaHashCheck confirma se a senha encriptada é correta
func SenhaHashCheck(SenhaHash, senha string) error {
	return bcrypt.CompareHashAndPassword([]byte(SenhaHash), []byte(senha))
}

// ENVIAR EMAIL #################

// Email ...
func Email(nome string) {
	from := mail.NewEmail("Rodrigo Valentergs", "valentergs@gmail.com")
	subject := "Bem-vindo " + nome + "!"
	to := mail.NewEmail(nome, "rodrigovalente@hotmail.com")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := nome + ", bem-vindo ao app"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
