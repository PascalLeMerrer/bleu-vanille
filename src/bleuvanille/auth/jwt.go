package auth

import (
	"bleuvanille/config"
	"bleuvanille/session"
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
)

// JWTAuth is a JSON Web Token middleware
func JWTAuth() echo.HandlerFunc {
	return func(context *echo.Context) error {
		// Skip WebSocket
		if (context.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
			return nil
		}

		header := context.Request().Header.Get("Authorization")
		if header == "" {
			cookie, err := context.Request().Cookie("token")
			if err == nil {
				header = cookie.Value
			}
		}
		prefixLength := len(Bearer)
		if len(header) > prefixLength+1 && header[:prefixLength] == Bearer {
			encodedToken := header[prefixLength+1:]
			token, err := ExtractToken(encodedToken)
			if err == nil {
				if token.Valid {
					// Store token data (=claims) in echo.Context
					// an aotuehticated user token contains a session Id
					context.Set("sessionId", token.Claims["id"])
					// a password reset token will contain an email
					context.Set("email", token.Claims["email"])
					return nil
				}
				deleteExpiredSession(token)
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
		}
		log.Printf("DEBUG: Invalid token. Header is %v\n", header)
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}

func deleteExpiredSession(token *jwt.Token) {
	if sessionID, ok := token.Claims["id"].(string); ok {
		session.Delete(sessionID)
	}
}

// ExtractToken decodes the token from a signed string representing the encoded token
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

// GetSessionToken generates a valid JWT token, and a unique session Id
func GetSessionToken() (string, string) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["typ"] = "JWT"
	token.Claims["exp"] = time.Now().Add(time.Hour * config.SessionTokenDuration).Unix()
	sessionID := uuid.NewV4().String()
	token.Claims["id"] = sessionID
	encodedToken, _ := token.SignedString([]byte(SigningKey))
	return encodedToken, sessionID
}

// GetResetToken generates a short lived JWT token, including the given email address
func GetResetToken(email string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["typ"] = "JWT"
	token.Claims["exp"] = time.Now().Add(time.Minute * config.ResetTokenDuration).Unix()
	token.Claims["email"] = email
	encodedToken, _ := token.SignedString([]byte(SigningKey))
	return encodedToken
}
