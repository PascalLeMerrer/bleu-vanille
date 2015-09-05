package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
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
	return func(c *echo.Context) error {

		// Skip WebSocket
		if (c.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
			return nil
		}

		auth := c.Request().Header.Get("Authorization")
		l := len(Bearer)
		httpError := echo.NewHTTPError(http.StatusUnauthorized)

		if len(auth) > l+1 && auth[:l] == Bearer {
			t, err := ExtractToken(auth[l+1:])
			if err == nil && t.Valid {
				// Store token claims in echo.Context
				c.Set("claims", t.Claims)
				return nil
			}
		}
		return httpError
	}
}

//ExtractToken decodes the token from a signed string representing the encoded token
func ExtractToken(signedString string) (*jwt.Token, error) {
	return jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {

		// Always check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the key for validation
		return []byte(SigningKey), nil
	})
}

//GetEncodedToken create a valid JWT token
func GetEncodedToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["typ"] = "JWT"
	token.Claims["exp"] = time.Now().Add(time.Hour * TokenDuration).Unix()
	encodedToken, _ := token.SignedString([]byte(SigningKey))
	return encodedToken
}
