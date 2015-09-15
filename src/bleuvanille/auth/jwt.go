package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/twinj/uuid"
)

const (
	// Bearer is the name of the header
	Bearer = "Bearer"
	// SigningKey is a secret
	SigningKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

	// TokenDuration defines the time to live of the token, in hours
	TokenDuration = 96
)

// JWTAuth is a JSON Web Token middleware
func JWTAuth() echo.HandlerFunc {
	return func(context *echo.Context) error {
		// Skip WebSocket
		if (context.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
			return nil
		}

		header := context.Request().Header.Get("Authorization")

		prefixLength := len(Bearer)
		httpError := echo.NewHTTPError(http.StatusUnauthorized)
		if len(header) > prefixLength+1 && header[:prefixLength] == Bearer {
			encodedToken := header[prefixLength+1:]
			token, err := ExtractToken(encodedToken)
			if err == nil && token.Valid {
				// Store token data (=claims) in echo.Context
				context.Set("sessionId", token.Claims["id"])
				return nil
			}
		}
		log.Printf("DEBUG: Invalid token. Header is %v\n", header)
		return httpError
	}
}

//ExtractToken decodes the token from a signed string representing the encoded token
func ExtractToken(signedString string) (*jwt.Token, error) {
	return jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {

		// Always check the signing method, otherwise there is a possible exploit
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the key for validation
		return []byte(SigningKey), nil
	})
}

//GetEncodedToken generates a valid JWT token, and a unique session Id
func GetEncodedToken() (string, string) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["typ"] = "JWT"
	token.Claims["exp"] = time.Now().Add(time.Hour * TokenDuration).Unix()
	sessionID := uuid.NewV4().String()
	token.Claims["id"] = sessionID
	encodedToken, _ := token.SignedString([]byte(SigningKey))
	return encodedToken, sessionID
}
