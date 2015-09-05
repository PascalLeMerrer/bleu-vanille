package main

import (
	"bleuvanille/auth"
	"bleuvanille/config"
	"bleuvanille/contact"
	"bleuvanille/user"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	//"github.com/lib/pq"
)

//Port on which the server is listening to
const (
	Port = 4000
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

	// myHTTPErrorHandler := func(err error, c *echo.Context) {
	// 	code := http.StatusInternalServerError
	// 	msg := http.StatusText(code)
	// 	if he, ok := err.(*echo.HTTPError); ok {
	// 		code = he.Code()
	// 		msg = he.Error()
	// 	}
	// 	fmt.Println(err)
	// 	if echoServer.Debug() {
	// 		msg = err.Error()
	// 		fmt.Println(msg)
	// 	}
	// 	if !c.Response().Commited() {
	// 		http.Error(c.Response(), msg, code)
	// 	}
	// 	log.Println(err)
	// }
	//
	// echoServer.SetHTTPErrorHandler(myHTTPErrorHandler)

	echoServer.Static("/js/", "public/js")
	echoServer.Static("/css/", "public/css")
	echoServer.Static("/fonts/", "public/fonts")
	echoServer.Static("/img/", "public/img")

	echoServer.Get("/", LandingPage)
	echoServer.Post("/contacts", contact.Create)
	echoServer.Get("/contacts", contact.Get)
	echoServer.Delete("/contacts", contact.Remove)
	echoServer.Post("/users", user.Create)
	// echo does not accept Delete request with body so we use a Post instead
	echoServer.Post("/users/delete", user.Remove)
	echoServer.Post("/users/login", user.Login)
	echoServer.Put("/users/password", user.ChangePassword)

	// Admin routes are protected by JWT
	restrictedRoutes := echoServer.Group("/admin")
	restrictedRoutes.Use(auth.JWTAuth())
	restrictedRoutes.Get("/contacts", contact.GetAll)

	log.Printf("Server listening to HTTP requests on port %d", Port)

	echoServer.Run(":" + strconv.Itoa(Port))
}
