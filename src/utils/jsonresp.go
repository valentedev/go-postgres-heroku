package utils

import (
	"encoding/json"
	"net/http"
)

// RespostaComErro ...
func RespostaComErro(w http.ResponseWriter, status int, erro string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(erro)
	return
}
