package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/joho/godotenv"
)

var csrfMd func(http.Handler) http.Handler

func init() {
	godotenv.Load()
	authKey := os.Getenv("CSRF_TOKEN_AUTH_KEY")
	secureMode := false
	if gin.Mode() == "release" {
		secureMode = true
	}
	fmt.Println("AUTH KEY: ", authKey)
	csrfMd = csrf.Protect([]byte(authKey),
		csrf.MaxAge(0),
		csrf.Secure(secureMode),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			msg := fmt.Sprintf(`{"message": "%s"}`, csrf.FailureReason(r))
			w.Write([]byte(msg))
		})),
	)
}

func CSRF() gin.HandlerFunc {
	return adapter.Wrap(csrfMd)
}
