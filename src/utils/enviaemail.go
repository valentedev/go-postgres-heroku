package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EnviaEmail para verificação e troca de senha
func EnviaEmail(nome, email, codigo string) {

	from := mail.NewEmail("Rodrigo Valente", "valentergs@gmail.com")
	subject := "Troca de senha - Admin.app"
	to := mail.NewEmail(nome, email)
	plainTextContent := "and easy to do anywhere, even with Go"
	// htmlContent := `
	// Clique no link abaixo para solicitar troca de sua senha.
	// http://localhost:8080/api/emailconfirma/` + codigo
	htmlContent := "Esse é seu código de validação: " + codigo

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
