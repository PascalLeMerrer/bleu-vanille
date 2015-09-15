package session

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// Middleware retrieves the session for an authenticated user
// It also deletes session for expired token
func Middleware() echo.HandlerFunc {
	return func(context *echo.Context) error {
		sessionID := context.Get("sessionId").(string)
		session, error := GetByID(sessionID)
		if error != nil {
			fmt.Printf("DEBUG: Session %v not found", sessionID)
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		context.Set("session", session)
		return nil
	}
}

// AdminMiddleware is a middleware that checks if the current user is an admin
func AdminMiddleware() echo.HandlerFunc {
	return func(context *echo.Context) error {
		session := context.Get("session").(*Session)
		log.Printf("DEBUG Session %v \n", session)
		if session.IsAdmin {
			return nil
		}
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}
