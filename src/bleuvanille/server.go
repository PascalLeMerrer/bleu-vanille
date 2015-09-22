package main

import (
	"bleuvanille/auth"
	"bleuvanille/config"
	"bleuvanille/contact"
	"bleuvanille/session"
	"bleuvanille/user"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	//"github.com/lib/pq"
)

// LandingPage displays the landing page
func LandingPage(context *echo.Context) error {
	return context.Render(http.StatusOK, "index", nil)
}

// Render processes a template
// name is the file name, without its HTML extension
func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Template type contains a list of templates
type Template struct {
	templates *template.Template
}

func main() {

	config.DatabaseInit()
	user.CreateDefault()

	echoServer := echo.New()
	echoServer.SetDebug(true)
	echoServer.ColoredLog(true)
	echoServer.Use(middleware.Logger())
	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.Gzip())
	// precompile templates
	templateRenderer := &Template{
		templates: template.Must(template.ParseGlob("src/bleuvanille/templates/*.html")),
	}
	echoServer.SetRenderer(templateRenderer)

	declareStaticRoutes(echoServer)
	declarePublicRoutes(echoServer)
	declarePrivateRoutes(echoServer)
	declareAdminRoutes(echoServer)
	declareSpecialRoutes(echoServer)
	addErrorHandler(echoServer)

	log.Printf("Server listening to HTTP requests on port %d\n", config.HostPort)

	echoServer.Run(":" + strconv.Itoa(config.HostPort))
}

// static pages
func declareStaticRoutes(echoServer *echo.Echo) {
	echoServer.Static("/js/", "public/js")
	echoServer.Static("/css/", "public/css")
	echoServer.Static("/fonts/", "public/fonts")
	echoServer.Static("/img/", "public/img")
}

// public pages
func declarePublicRoutes(echoServer *echo.Echo) {
	echoServer.Get("/", LandingPage)
	echoServer.Post("/contacts", contact.Create)
	echoServer.Post("/users", user.Create)
	echoServer.Post("/users/login", user.Login)
	echoServer.Post("/users/sendresetlink", user.SendResetLink)
	echoServer.Get("/users/resetform", user.DisplayResetForm)
}

// privates Routes require a valid user auth token and a sessionID
func declarePrivateRoutes(echoServer *echo.Echo) {
	userRoutes := echoServer.Group("/users")
	userRoutes.Use(auth.JWTAuth())
	userRoutes.Use(session.Middleware())

	// echo does not accept Delete request with body so we use a Post instead
	userRoutes.Post("/delete", user.Remove)
	userRoutes.Put("/password", user.ChangePassword)
	userRoutes.Get("/:userID/profile", user.Profile)
}

// special Routes require a valid user auth token but no sessionID
func declareSpecialRoutes(echoServer *echo.Echo) {
	specialRoutes := echoServer.Group("/special")
	specialRoutes.Use(auth.JWTAuth())
	specialRoutes.Post("/resetpassword", user.ResetPassword)
}

// Admin routes require a valid auth token AND the user to have the admin rights
func declareAdminRoutes(echoServer *echo.Echo) {

	adminRoutes := echoServer.Group("/admin")
	adminRoutes.Use(auth.JWTAuth())
	adminRoutes.Use(session.Middleware())
	adminRoutes.Use(session.AdminMiddleware())
	adminRoutes.Get("/contacts", contact.GetAll)
	adminRoutes.Delete("/contacts", contact.Remove)
}

// Defines a custom error handler
// Is only invoked by echo when the error occurs in the handlerFunction
// and not when a middleware returns an error :(
// TODO: see how to improve this
func addErrorHandler(echoServer *echo.Echo) {

	myHTTPErrorHandler := func(err error, context *echo.Context) {
		fmt.Println("Custom error handler invoked")
		code := http.StatusInternalServerError
		message := http.StatusText(code)
		if httpError, ok := err.(*echo.HTTPError); ok {
			code = httpError.Code()
			message = httpError.Error()
		}
		fmt.Println(err)
		if echoServer.Debug() {
			message = err.Error()
			fmt.Println(message)
		}
		if !context.Response().Commited() {
			http.Error(context.Response(), message, code)
		}
		log.Println(err)
	}

	echoServer.SetHTTPErrorHandler(myHTTPErrorHandler)
}
