package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/valentedev/go-postgres-heroku/src/api"
	"github.com/valentedev/go-postgres-heroku/src/controllers.go"
	"github.com/valentedev/go-postgres-heroku/src/utils"
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

	mux.HandleFunc("/", controllers.Home())
	mux.HandleFunc("/admin/login/", controllers.Login(conect))
	mux.HandleFunc("/admin/login/redirect/", controllers.Logado(conect))
	mux.HandleFunc("/admin/", utils.TokenMiddleware(controllers.AdminHome(conect)))
	mux.HandleFunc("/admin/usuario/", utils.TokenMiddleware(controllers.Usuario(conect)))
	mux.HandleFunc("/admin/usuario/criar/", utils.TokenMiddleware(controllers.CriarUsuario()))
	mux.HandleFunc("/admin/usuario/criado/", utils.TokenMiddleware(controllers.Criado(conect)))
	mux.HandleFunc("/admin/usuario/editar/", utils.TokenMiddleware(controllers.EditarUsuario(conect)))
	mux.HandleFunc("/admin/usuario/editado/", utils.TokenMiddleware(controllers.Editado(conect)))
	mux.HandleFunc("/admin/usuario/deletar/", utils.TokenMiddleware(controllers.Deletar(conect)))
	mux.HandleFunc("/admin/usuario/deletado/", utils.TokenMiddleware(controllers.Deletado(conect)))
	mux.HandleFunc("/admin/usuario/novasenha/", utils.TokenMiddleware(controllers.NovaSenha(conect)))
	mux.HandleFunc("/admin/usuario/novasenha/confirma/", utils.TokenMiddleware(controllers.NovaSenhaConfirma(conect)))
	mux.HandleFunc("/api/", api.Home())
	mux.HandleFunc("/api/login", api.Login(conect))
	mux.HandleFunc("/api/cadastro", api.Cadastro(conect))
	mux.HandleFunc("/api/reset", api.ResetSenhaUm(conect))
	mux.HandleFunc("/api/reset/confirma", api.ResetSenhaDois(conect))
	mux.HandleFunc("/api/mudarsenha", api.MudarSenha(conect))
	//mux.HandleFunc("/api/emailconfirma", api.EmailConfirma(conect))

	// // aqui chamamos a func seed() para migrar os dados do []UsuariosDB para Banco de Dados novo.
	// // depois que os dados foram migrados, podem deixar de chamar a função seed(db *sql.DB)
	//utils.seed(conect)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://valentedev-react-auth.herokuapp.com"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		Debug:          true,
	})

	handler := c.Handler(mux)

	addr := ":" + port
	log.Println("Listen on port 8080...")
	err = http.ListenAndServe(addr, handler)
	log.Fatal(err)

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
