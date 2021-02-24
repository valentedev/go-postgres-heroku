package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/valentedev/go-postgres-heroku/src/usuarios"
	"github.com/valentedev/go-postgres-heroku/src/utils"
)

// Home - home para os endpoints
func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		utils.Logging(r)

		var tpl *template.Template

		tpl = template.Must(template.ParseGlob("./templates/*"))

		err := tpl.ExecuteTemplate(w, "API", nil)
		if err != nil {
			panic(err)
		}
	}
}

// Login recebe um JSON com email+senha, consulta o DB e retorno um token de acesso
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var usuario usuarios.Usuarios

		if r.Method != "POST" {
			http.Error(w, "Método não autorizado", 405)
			return
		}

		utils.Logging(r)

		json.NewDecoder(r.Body).Decode(&usuario)
		senhaJSON := usuario.Senha

		query := `SELECT id, nome, sobrenome, email, senha, admin, ativo FROM usuarios WHERE email=$1;`
		row := db.QueryRow(query, usuario.Email)
		err := row.Scan(&usuario.ID, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, &usuario.Senha, &usuario.Admin, &usuario.Ativo)
		if err != nil {
			utils.RespostaComErro(w, 404, "Usuário não encontrado")
			return
		}

		t, err := utils.TokenAPI(usuario)
		if err != nil {
			http.Error(w, "Não foi possivel retornar un Token", 404)
			return
		}

		if utils.SenhaHashCheck(usuario.Senha, senhaJSON) == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
		} else {
			//http.Error(w, "Senha inválida", 401)
			utils.RespostaComErro(w, 401, "Senha inválida")
			return
		}
	}
}

// Cadastro recebe un JSON com dados de nome, sobrenome, email e senha e retorna um mensagem de sucesso
func Cadastro(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var usuario usuarios.Usuarios

		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
			return

		}

		utils.Logging(r)

		json.NewDecoder(r.Body).Decode(&usuario)

		senhaJSON := usuario.Senha
		senhaEncrip, err := utils.SenhaHash(senhaJSON)
		if err != nil {
			panic(err)
		}

		query := `INSERT INTO usuarios (nome, sobrenome, email, senha, admin, ativo) VALUES ($1,$2,$3,$4,$5,$6);`
		_, err = db.Exec(query, &usuario.Nome, &usuario.Sobrenome, &usuario.Email, senhaEncrip, "false", "false")
		if err != nil {
			panic(err)
		}

		sucesso := "Usuário cadastrado com sucesso!"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sucesso)

	}

}

// MudarSenhaStruct Recebe JSON com informações de usuário que quer mudar senha
type MudarSenhaStruct struct {
	Token string `json:"token"`
	Senha string `json:"senha"`
}

// MudarSenha recebe um JWT e nova senha e
func MudarSenha(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mudarSenha MudarSenhaStruct

		if r.Method != "POST" {
			//http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
			utils.RespostaComErro(w, 405, "Método não autorizado")
			return
		}

		utils.Logging(r)
		json.NewDecoder(r.Body).Decode(&mudarSenha)
		senha, err := utils.SenhaHash(mudarSenha.Senha)
		if err != nil {
			utils.RespostaComErro(w, 401, "Aconteceu algum problema. Tente novamente.")
			return
		}

		token := mudarSenha.Token

		err = utils.TokenCheck(token)
		if err != nil {
			utils.RespostaComErro(w, 401, "Token inválido")
			return
		}

		email := utils.TokenAPIEmail(token)

		query := `UPDATE usuarios SET senha=$1 WHERE email=$2;`
		_, err = db.Exec(query, senha, email)
		if err != nil {
			utils.RespostaComErro(w, 404, "Usuário não encontrado")
			return
		}

		mensagem := "Senha alterada com sucesso"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mensagem)

	}
}
