package api

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/valentedev/go-postgres-heroku/src/usuarios"
	"github.com/valentedev/go-postgres-heroku/src/utils"
)

// Vercod é uma struct para verificação de codigo
type Vercod struct {
	ID      int
	Criacao string
	Usuario int
	Codigo  string
}

// ResetSenha ...
func ResetSenha(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var usuario usuarios.Usuarios

		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
			return

		}

		utils.Logging(r)

		json.NewDecoder(r.Body).Decode(&usuario)

		query := `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE email=$1;`
		row := db.QueryRow(query, usuario.Email)
		err := row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			utils.RespostaComErro(w, 404, "Usuário não encontrado")
			return
		}

		codigo := CodigoVerificação(16)

		id := usuario.ID
		nome := usuario.Nome
		email := usuario.Email

		query = `INSERT INTO vercod (usuario, codigo) VALUES ($1,$2)`
		_, err = db.Exec(query, id, codigo)
		if err != nil {
			panic(err)
		}

		utils.EnviaEmail(nome, email, codigo)
	}
}

// CodigoVerificação é um código aleatório de 8 digitos que é enviado ao usuário para verificação autenticidade
func CodigoVerificação(n int) string {
	const alfaBeta = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		x := rand.Int63() % int64(len(alfaBeta))
		sb.WriteByte(alfaBeta[x])
	}
	return sb.String()
}

// EmailConfirma recebe um link com o codigo de verificação
func EmailConfirma(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var usuario usuarios.Usuarios
		var vercod Vercod

		if r.Method != "GET" {
			http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
			return

		}

		utils.Logging(r)

		// id será o código de verificação informado no URL
		params := r.URL.Path
		id := strings.TrimPrefix(params, "/api/emailconfirma/")

		// faz uma consulta no BD se o id (código de verificação) existe
		query := `SELECT id, criado_em, usuario, codigo FROM vercod WHERE codigo=$1;`
		row := db.QueryRow(query, id)
		// coloca o resultado da consulta no struct Vercod
		err := row.Scan(&vercod.ID, &vercod.Criacao, &vercod.Usuario, &vercod.Codigo)
		if err != nil {
			panic(err)
		}

		// aloca a data de criação encontrada na variável criacaoVercod
		criacaoVercod := vercod.Criacao
		// formata criacaoVercod para time.Time, formato RFC3339
		inicio, err := time.Parse(time.RFC3339, criacaoVercod)
		if err != nil {
			panic(err)
		}

		fim := time.Now()
		// estabelece a diferença de tempo entre a criação do código de verificação e o momento da consulta
		delta := fim.Sub(inicio)

		// se delta for maior que 10 min retorna JSON com mensagem
		if delta > (time.Minute * 10) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Seu código de verificação está vencido")
			return
		}

		// consulta o BD para trazer o usuário informado no vercod
		query = `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE id=$1;`
		row = db.QueryRow(query, vercod.Usuario)
		err = row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			panic(err)
		}

		// emite um token com esse o usuário
		// TODO: criar TokenAPI que será usado apenas pelo usuário, não Admin
		token, err := utils.Token(usuario)
		if err != nil {
			panic(err)
		}

		// responde com um JSON + token com usuário. Esse usuário será comparado com o usuário logado no Frontend.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)

	}

}
