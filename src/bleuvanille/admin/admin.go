package admin

import (
	"net/http"

	"github.com/labstack/echo"
)

// LoginPage returns the admin login page
func LoginPage(context *echo.Context) error {
	// TODO: redirect to dashboard when the user is logged with admin account
	return context.Render(http.StatusOK, "admin/login", nil)
}
