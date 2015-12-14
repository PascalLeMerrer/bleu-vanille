package admin

import (
	"bleuvanille/session"
	"net/http"

	"github.com/labstack/echo"
)

// LoginPage returns the admin login page
func LoginPage(context *echo.Context) error {
	return context.Render(http.StatusOK, "admin/login", nil)
}

// Dashboard displays the main administration page
func Dashboard(context *echo.Context) error {
	session := context.Get("session").(*session.Session)
	return context.Render(http.StatusOK, "admin/dashboard", session)
}