package helper

import (
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	UserRole  string
}

var JWT_SECRET_KEY string = os.Getenv("JWT_SECRET_KEY")
var secretKey = []byte(JWT_SECRET_KEY)

// CreateToken creates a signed JWT string for the given email that expires in 24 hours
func CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": email,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verifies the authenticity of the given token
//
// If no error occurs, then the token is valid.
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

type PasswordComplexityResults struct {
	// Whether or not the password meets the criteria for expected number of upper-case characters
	HasExpectedUpperCaseCt bool
	// Whether or not the password meets the criteria for expected number of special characters
	HasExpectedSpecialCaseCt bool
	// Whether or not the password meets the criteria for expected number of digit characters
	HasExpectedDigitCt bool
	// Whether or not the password meets the criteria for expected minimum length
	HasMinLength bool
}

func VerifyPasswordComplexity(password string, digitCt int, specialCt int, upperCaseCt int, minLength int) *PasswordComplexityResults {
	var result PasswordComplexityResults
	capitals := 0
	digits := 0
	specials := 0
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			digits++
		case unicode.IsUpper(c):
			capitals++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			specials++
		}
	}
	result.HasExpectedDigitCt = digits >= digitCt
	result.HasExpectedSpecialCaseCt = specials >= specialCt
	result.HasExpectedUpperCaseCt = capitals >= upperCaseCt
	result.HasMinLength = len(password) >= minLength
	return &result
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
