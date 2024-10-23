package auth

import (
	"net/http"
)

type Authenticator interface {
	PerformAuthentication(writer http.ResponseWriter, request *http.Request) (string, http.ResponseWriter)
}
