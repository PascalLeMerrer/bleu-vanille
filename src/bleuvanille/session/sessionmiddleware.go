package session

import (
	"net/http"

	"github.com/labstack/echo"
)

// Middleware retrieves the session for an authenticated user
// It also deletes session for expired token
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		rawSessionID := context.Get("sessionId")
		if rawSessionID == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		sessionID := rawSessionID.(string)
		session, error := FindByID(sessionID)

		if error != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		context.Set("session", session)
		return next(context)
	}
}

// AdminMiddleware is a middleware that checks if the current user is an admin
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		session := context.Get("session").(*Session)

		if session != nil && session.IsAdmin {
			return next(context)
		}
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}
