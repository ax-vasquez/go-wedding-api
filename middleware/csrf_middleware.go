package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

var csrfMd func(http.Handler) http.Handler

func init() {
	secureMode := false
	if gin.Mode() == "release" {
		secureMode = true
	}
	authKey := os.Getenv("CSRF_TOKEN_AUTH_KEY")
	csrfMd = csrf.Protect([]byte(authKey),
		csrf.MaxAge(0),
		csrf.Secure(secureMode),
		csrf.TrustedOrigins([]string{"*"}),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message": "Forbidden - CSRF token invalid"}`))
		})),
	)
}

func CSRF() gin.HandlerFunc {
	return adapter.Wrap(csrfMd)
}
