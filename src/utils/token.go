package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/valentedev/go-postgres-heroku/src/usuarios"
)

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
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			Issuer:    "go-postgres-heroku",
		},
		Email: u.Email,
		Nome:  u.Nome,
	}

	// evita que alguem logado com não-admin use esse token
	if u.Admin == false {
		return "", fmt.Errorf("Não autorizado %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	tokenString, err := token.SignedString([]byte(assinatura))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

// TokenAPI envia uma string para o client que será usada para autenticação.
func TokenAPI(u usuarios.Usuarios) (string, error) {

	var err error

	claims := minhasClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
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

// TokenAPIEmail recebe um JWT e retorna uma string com o e-mail do usuário
func TokenAPIEmail(s string) string {

	afterVerificationToken, err := jwt.ParseWithClaims(s, &minhasClaims{}, func(beforeVeritificationToken *jwt.Token) (interface{}, error) {
		return []byte(assinatura), nil
	})
	if err != nil {
		panic(err)
	}

	tokenOK := afterVerificationToken.Valid && err == nil

	claims := afterVerificationToken.Claims.(*minhasClaims)

	var email string
	email = claims.Email

	if !tokenOK {
		return "Token inválido"
	}

	return email
}

// TokenCheck verifica se o Token é válido
func TokenCheck(t string) error {
	afterVerificationToken, err := jwt.ParseWithClaims(t, &minhasClaims{}, func(beforeVeritificationToken *jwt.Token) (interface{}, error) {
		return []byte(assinatura), nil
	})
	if err != nil || afterVerificationToken.Valid == false {
		return err
	}

	return nil
}

// Payload armazena os dados retirados de um token válido
type Payload struct {
	Nome  string
	Email string
}

// TokenPayload verifica a validade do token e retorna um struct com dados do token.payload
func TokenPayload(c *http.Cookie) Payload {

	tokenString := c.Value

	afterVerificationToken, err := jwt.ParseWithClaims(tokenString, &minhasClaims{}, func(beforeVeritificationToken *jwt.Token) (interface{}, error) {
		return []byte(assinatura), nil
	})
	if err != nil {
		panic(err)
	}

	tokenOK := afterVerificationToken.Valid && err == nil

	claims := afterVerificationToken.Claims.(*minhasClaims)

	var payload Payload

	if tokenOK {
		payload.Nome = claims.Nome
		payload.Email = claims.Email
		return payload
	}

	return payload

}

// TokenMiddleware é um wrapper que vai verificar se há um token válido em cada Handler.
func TokenMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/admin/login/", http.StatusSeeOther)
			return
		}

		tokenString := c.Value

		tokenVerificado, err := jwt.ParseWithClaims(tokenString, &minhasClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(assinatura), nil
		})
		if err != nil || !tokenVerificado.Valid {
			fmt.Println("Token inválido ou inexistente")
			http.Redirect(w, r, "/admin/login/", 307)
			return
		}

		next.ServeHTTP(w, r)
	})

}
