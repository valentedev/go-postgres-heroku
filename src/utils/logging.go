package utils

import (
	"fmt"
	"net/http"
	"time"
)

// Logging ...
func Logging(r *http.Request) {
	metodo := r.Method
	params := r.URL.Path
	host := r.Host
	fmt.Printf("Request: %s %s%s - %s\n", metodo, host, params, time.Now().Format(time.RFC3339))
}
