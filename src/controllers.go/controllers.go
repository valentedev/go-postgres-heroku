package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/valentedev/go-postgres-heroku/src/usuarios"
	"github.com/valentedev/go-postgres-heroku/src/utils"
)

// Home é a página de inicio do aplicativo
func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		utils.Logging(r)

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		var tpl *template.Template

		tpl = template.Must(template.ParseGlob("./templates/*"))

		err := tpl.ExecuteTemplate(w, "Index", nil)
		if err != nil {
			panic(err)
		}
	}
}

// AdminHome é uma função que vai usar o Template index.html e injeta informações de usuarios em uma tabela
func AdminHome(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		utils.Logging(r)

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

		c, err := r.Cookie("session")
		if err != nil {
			c = &http.Cookie{}
		}

		t := utils.TokenPayload(c)

		type Dados struct {
			Linhas  []usuarios.Usuarios
			Usuario string
		}

		dados := Dados{
			Linhas:  linhas,
			Usuario: t.Nome,
		}

		var tpl *template.Template

		tpl = template.Must(template.ParseGlob("./templates/*"))

		err = tpl.ExecuteTemplate(w, "Admin", dados)
		if err != nil {
			panic(err)
		}
	}

}

//Usuario vai retornar um usuario que tenha o mesmo ID informado no http request
func Usuario(db *sql.DB) http.HandlerFunc {

	//retornamos um Handler será uma função anônima
	return func(w http.ResponseWriter, r *http.Request) {

		utils.Logging(r)

		//"params" é o URL de request. Nesse caso, /usuario/{id}
		params := r.URL.Path
		//"id" é o params sem /usuario/. Ficamos apenas com o numero que nos interessa: {id}
		id := strings.TrimPrefix(params, "/admin/usuario/")
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

		c, err := r.Cookie("session")
		if err != nil {
			c = &http.Cookie{}
		}

		tokenEmail := utils.TokenPayload(c)
		if err != nil {
			panic(err)
		}

		type Dados struct {
			Usuario    usuarios.Usuarios
			TokenEmail string
		}

		dados := Dados{
			Usuario:    usuario,
			TokenEmail: tokenEmail.Email,
		}

		//Criamos um template tpl
		tpl := template.Must(template.ParseGlob("./templates/*"))
		//executamos o template com os dados presentes em "usuario" e enviamos o "response w"
		err = tpl.ExecuteTemplate(w, "Detalhes", dados)
		if err != nil {
			panic(err)
		}
	}
}

//CriarUsuario gera um formulário para entrada de dados no DB
func CriarUsuario() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Logging(r)

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

		utils.Logging(r)

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
			http.Redirect(w, r, "/admin/usuario/criar/", http.StatusSeeOther)
		}

		senhaByte, err := utils.SenhaHash(senha)
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

		utils.Logging(r)

		//"params" é o URL de request. Nesse caso, /usuario/{id}
		params := r.URL.Path
		//"id" é o params sem /usuario/. Ficamos apenas com o numero que nos interessa: {id}
		id := strings.TrimPrefix(params, "/admin/usuario/editar/")
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

		utils.Logging(r)

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

		utils.Logging(r)

		params := r.URL.Path
		id := strings.TrimPrefix(params, "/admin/usuario/deletar/")
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
		utils.Logging(r)

		params := r.URL.Path
		id := strings.TrimPrefix(params, "/admin/usuario/deletado/")
		idint, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("invalid param format")
		}

		query := `DELETE FROM usuarios WHERE id=$1;`
		_, err = db.Exec(query, idint)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/admin/", 307)

	}
}

//NovaSenha é um handler para mudar a senha do admin logado
func NovaSenha(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		utils.Logging(r)

		c, err := r.Cookie("session")
		if err != nil {
			c = &http.Cookie{}
		}

		tokenEmail := utils.TokenPayload(c).Email
		if err != nil {
			panic(err)
		}

		//"params" é o URL de request. Nesse caso, /usuario/{id}
		params := r.URL.Path
		//"id" é o params sem /usuario/. Ficamos apenas com o numero que nos interessa: {id}
		id := strings.TrimPrefix(params, "/admin/usuario/novasenha/")
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

		// Caso o email do usuario seja diferente do email do token, o processo é interrompido. Assim se evita que um admin altere a senha outro usuário.
		if tokenEmail != usuario.Email {
			http.Error(w, "Unauthorized", 401)
			return
		}

		var tpl *template.Template
		tpl = template.Must(template.ParseGlob("./templates/*"))
		err = tpl.ExecuteTemplate(w, "NovaSenha", usuario)
		if err != nil {
			panic(err)
		}
	}

}

// NovaSenhaConfirma confirma que uma nova senha foi criada
func NovaSenhaConfirma(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		idint, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		senha := r.FormValue("senha")
		senhaByte, err := utils.SenhaHash(senha)
		if err != nil {
			fmt.Println(err)
		}
		senha = string(senhaByte)
		var usuario usuarios.Usuarios
		query := `UPDATE usuarios SET senha=$1 WHERE id=$2;`
		_, err = db.Exec(query, senha, idint)
		if err != nil {
			panic(err)
		}
		query = `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE id=$1;`
		row := db.QueryRow(query, idint)
		err = row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			fmt.Println(err)
		}

		tpl := template.Must(template.ParseGlob("./templates/*"))
		err = tpl.ExecuteTemplate(w, "NovaSenhaCriada", usuario)
		if err != nil {
			panic(err)
		}
	}
}

// Login recebe email e senha e autentica acesso
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		utils.Logging(r)

		c := http.Cookie{
			Path:     "/",
			Name:     "session",
			Value:    "",
			HttpOnly: true,
			Expires:  time.Unix(0, 0),
			Secure:   false,
		}

		http.SetCookie(w, &c)

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

		utils.Logging(r)

		if r.Method != http.MethodPost {
			http.Error(w, "Método não autorizado", 405)
			return
		}

		email := r.FormValue("email")
		senha := r.FormValue("senha")
		if email == "" || senha == "" {
			http.Redirect(w, r, "/admin/login/", http.StatusSeeOther)
		}

		// query baseada em email
		query := `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE email=$1`
		row := db.QueryRow(query, email)
		var usuario usuarios.Usuarios
		err := row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			fmt.Println(err)
		}

		if utils.SenhaHashCheck(usuario.Senha, senha) != nil || usuario.Admin != true {
			http.Error(w, "Acesso não autorizado", 401)
		}

		token, err := utils.Token(usuario)
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

		http.SetCookie(w, &c)
		//w.Header().Add("Authorization", token)
		//Email(nome)
		http.Redirect(w, r, "/admin/", 307)

	}
}
