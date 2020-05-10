package middleware

import (
	"net/http"
)

func JsonContentTypeMiddleware(next http.Handler) http.Handler {
	return CreateHeadersMiddleware(map[string]string{
		"Content-Type": "application/json",
	})(next)
}
