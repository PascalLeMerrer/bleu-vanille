package bleuvanille

import (
	"bleuvanille/admin"
	"bleuvanille/auth"
	"bleuvanille/config"
	"bleuvanille/contact"
	"bleuvanille/ingredient"
	"bleuvanille/log"
	"bleuvanille/session"
	"bleuvanille/statistics"
	"bleuvanille/user"
	"fmt"
	gommonlog "github.com/labstack/gommon/log"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

// the id of the git commit
var Sha1 string

// Render processes a template
// name is the file name, without its HTML extension
func (t *Template) Render(w io.Writer, name string, data interface{}, context echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Template type contains a list of templates
type Template struct {
	templates *template.Template
}

func main() {

	// dependencies injection

	user.CreateDefault()

	echoServer := echo.New()
	if config.ServerDebug() {
		echoServer.SetDebug(true)
		echoServer.SetLogLevel(gommonlog.DEBUG)
	}

	echoServer.Use(middleware.Logger())
	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.Gzip())
	echoServer.Use(log.Middleware)

	// precompile templates

	templates := template.New("template")
	templates = template.Must(template.ParseGlob("src/bleuvanille/templates/*.html")) // parse templates in root dir
	// Parse templates in subdir
	filepath.Walk("src/bleuvanille/templates", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			_, err := templates.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	templateRenderer := &Template{
		templates: templates,
	}

	echoServer.SetRenderer(templateRenderer)

	declareStaticRoutes(echoServer)
	declarePublicRoutes(echoServer)
	declarePrivateRoutes(echoServer)
	declareIngredientRoutes(echoServer)
	declareAdminRoutes(echoServer)
	declareSpecialRoutes(echoServer)
	addErrorHandler(echoServer)

	fmt.Printf("Server listening to HTTP requests on port %d\n", config.HostPort)

	echoServer.Run(standard.New("127.0.0.1:" + strconv.Itoa(config.HostPort)))
}

// static pages
func declareStaticRoutes(echoServer *echo.Echo) {
	echoServer.Static("/js/", "public/js")
	echoServer.Static("/css/", "public/css")
	echoServer.Static("/fonts/", "public/fonts")
	echoServer.Static("/img/", "public/img")
	echoServer.Static("/tags/", "public/tags")
	echoServer.Static("/robots.txt", "public/robots.txt")
}

// public pages
func declarePublicRoutes(echoServer *echo.Echo) {
	echoServer.Get("/", contact.LandingPage())
	echoServer.Get("/version", getVersion())
	echoServer.Get("/admin", admin.LoginPage())
	echoServer.Get("/admin/", admin.LoginPage())
	echoServer.Post("/contacts", contact.Create())
	echoServer.Post("/users", user.Create())
	echoServer.Post("/users/login", user.Login())
	echoServer.Post("/users/sendResetLink", user.SendResetLink())
	echoServer.Get("/users/resetForm", user.DisplayResetForm())
}

// privates Routes require a valid user auth token and a sessionID
func declarePrivateRoutes(echoServer *echo.Echo) {
	userRoutes := echoServer.Group("/users")
	userRoutes.Use(auth.JWTAuth)
	userRoutes.Use(session.Middleware)

	userRoutes.Post("/logout", session.Logout())
	// echo does not accept Delete request with body so we use a Post instead
	userRoutes.Post("/delete", user.Delete())
	userRoutes.Put("/password", user.ChangePassword())
	userRoutes.Get("/:userID", user.Profile())
	userRoutes.Patch("/:userID", user.Patch())
}

// ingredients Routes require a valid user auth token and a sessionID
func declareIngredientRoutes(echoServer *echo.Echo) {
	ingredientRoutes := echoServer.Group("/ingredients")
	ingredientRoutes.Use(auth.JWTAuth)
	ingredientRoutes.Use(session.Middleware)
	ingredientRoutes.Get("/:key", ingredient.Get())
	ingredientRoutes.Get("", ingredient.GetAll())
	ingredientRoutes.Post("", ingredient.Create())
	ingredientRoutes.Patch("/:key", ingredient.Patch())
	ingredientRoutes.Delete("/:key", ingredient.Delete())
}

// special Routes require a valid user auth token but no sessionID
func declareSpecialRoutes(echoServer *echo.Echo) {
	specialRoutes := echoServer.Group("/special")
	specialRoutes.Use(auth.JWTAuth)
	specialRoutes.Post("/resetPassword", user.ResetPassword())
}

// Admin routes require a valid auth token AND the user to have the admin rights
func declareAdminRoutes(echoServer *echo.Echo) {

	adminRoutes := echoServer.Group("/admin")
	adminRoutes.Use(auth.JWTAuth)
	adminRoutes.Use(session.Middleware)
	adminRoutes.Use(session.AdminMiddleware)

	adminRoutes.Get("/dashboard", admin.Dashboard())
	adminRoutes.Get("/contacts", contact.GetAll())
	adminRoutes.Get("/users", user.GetAll())
	adminRoutes.Get("/users/email", user.Get())
	adminRoutes.Delete("/users/:userID", user.RemoveByAdmin())
	adminRoutes.Delete("/contacts", contact.Delete())
	adminRoutes.Get("/statistics/:counter", statistics.Get())
}

func getVersion() echo.HandlerFunc {
	return func(context echo.Context) error {
		version := struct {
			Version string `json:"version"`
		}{Sha1}
		return context.JSON(200, version)
	}
}

// Defines a custom error handler
// Is only invoked by echo when the error occurs in the handlerFunction
// and not when a middleware returns an error :(
// TODO: see how to improve this
func addErrorHandler(echoServer *echo.Echo) {

	myHTTPErrorHandler := func(err error, context echo.Context) {
		fmt.Printf("Custom error handler invoked: %s\n", err)
		code := http.StatusInternalServerError
		message := http.StatusText(code)
		if httpError, ok := err.(*echo.HTTPError); ok {
			code = httpError.Code
			message = httpError.Error()
		}

		log.Error(context, err.Error())

		if !context.Response().Committed() {
			http.Error(context.Response().(*standard.Response).ResponseWriter, message, code)
		}
	}

	echoServer.SetHTTPErrorHandler(myHTTPErrorHandler)
}
