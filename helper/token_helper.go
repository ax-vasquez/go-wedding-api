package helper

import (
	"log"
	"os"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	*jwt.StandardClaims
	TokenType string
	ID        uuid.UUID `json:"id"`
	// The user's role, which can be "GUEST", "INVITEE" or "ADMIN". Defaults to "GUEST".
	Role string `json:"role"`
	// The user's first name.
	FirstName string `json:"first_name"`
	// The user's last name.
	LastName string `json:"last_name"`
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

// GenerateAllTokens generates a signed token and signed refresh token
func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid uuid.UUID) (signedToken string, signedRefreshToken string, err error) {

	// Claims to be stored in the token
	claims := &CustomClaims{
		ID:        uid,
		FirstName: firstName,
		LastName:  lastName,
		Role:      userType,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	// Claims to be stored in the refresh token
	refreshClaims := &CustomClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secretKey)
	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (claims *CustomClaims, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil

		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		msg = "the token is invalid"
		return nil, msg
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return nil, msg
	}

	return claims, msg
}

func VerifyPassword(userPassword string, providedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword)); err != nil {
		return false
	}
	return true
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
