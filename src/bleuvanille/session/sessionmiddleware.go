package session

import (
	"bleuvanille/log"
	"fmt"
	"net/http"
	"github.com/labstack/echo"
)

// Middleware retrieves the session for an authenticated user
// It also deletes session for expired token
func Middleware() echo.HandlerFunc {
	return func(context *echo.Context) error {
		rawSessionID := context.Get("sessionId")
		if rawSessionID == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		sessionID := rawSessionID.(string)
		session, error := GetByID(sessionID)
		
		log.Debug(context, "SessionId generated : " + sessionID)
		
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
		
		log.Debug(context, fmt.Sprintf("Session %v \n", session))
		if session.IsAdmin {
			return nil
		}
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}
