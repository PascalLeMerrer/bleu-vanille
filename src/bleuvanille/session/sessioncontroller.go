package session

import (
	"net/http"

	"github.com/labstack/echo"
)

// Logout deletes the session for the current user, and so invalidate its Authorization header
func Logout() echo.HandlerFunc {
	return func(context echo.Context) error {
		rawSessionID := context.Get("sessionId")
		if rawSessionID == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		sessionID := rawSessionID.(string)
		err := Remove(sessionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return context.JSON(http.StatusOK, "")
	}
}
