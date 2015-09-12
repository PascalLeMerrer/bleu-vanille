package user

import (
	"bleuvanille/auth"
	"bleuvanille/session"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
)

// Create creates a new user
func Create(context *echo.Context) error {

	email := context.Form("email")
	if email == "" {
		return context.JSON(http.StatusBadRequest, errors.New("Missing email parameter in POST body"))
	}
	firstname := context.Form("firstname")
	if firstname == "" {
		return context.JSON(http.StatusBadRequest, errors.New("Missing firstname parameter in POST body"))
	}
	lastname := context.Form("lastname")
	if lastname == "" {
		return context.JSON(http.StatusBadRequest, errors.New("Missing lastname parameter in POST body"))
	}
	password := context.Form("password")
	if password == "" {
		return context.JSON(http.StatusBadRequest, errors.New("Missing password parameter in POST body"))
	}

	user, err := NewUser(email, firstname, lastname, password)

	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusInternalServerError, errors.New("User creation error"))
	}

	_, err = SaveUser(user)
	if err != nil {
		//TODO: see if an error code is available and could be used
		if err, ok := err.(*pq.Error); ok {
			// see http://www.postgresql.org/docs/9.3/static/errcodes-appendix.html
			if err.Code.Name() == "unique_violation" {
				return context.JSON(http.StatusConflict, errors.New("User already exists"))
			}
		}

		return context.JSON(http.StatusInternalServerError, errors.New("User account creation failed - cannot be saved"))
	}
	user.Hash = "" // dont return the hash, for security concerns

	return context.JSON(http.StatusCreated, user)
}

// Get returns the entry for a given email if any
func Get(context *echo.Context) error {

	email := context.Query("email")

	if email == "" {
		return context.JSON(http.StatusBadRequest, errors.New("Missing email parameter in GET request"))
	}

	var user, err = LoadByEmail(email)
	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusNotFound, nil)
	}
	return context.JSON(http.StatusOK, user)
}

// GetAll writes the list of all users
func GetAll(context *echo.Context) error {
	users, err := LoadAll()
	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusInternalServerError, errors.New("User list retrieval error"))
	}
	return context.JSON(http.StatusOK, users)
}

// Remove removes the user account for a given email
func Remove(context *echo.Context) error {

	//TODO many lines to be factorised with login
	// try to use this pattern: http://talks.golang.org/2013/go4python.slide#37
	email, password := context.Form("email"), context.Form("password")
	if email == "" {
		log.Println("Missing email parameter in DELETE request")
		return context.JSON(http.StatusBadRequest, errors.New("Missing email parameter in DELETE request"))
	}
	if password == "" {
		log.Println("Missing password parameter in DELETE request")
		return context.JSON(http.StatusBadRequest, errors.New("Missing password parameter in DELETE request"))
	}
	user, err := LoadByEmail(email)

	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusUnauthorized, errors.New("Wrong email and password combination"))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		return context.JSON(http.StatusUnauthorized, errors.New("Wrong email and password combination"))
	}

	deleteErr := Delete(user)
	if deleteErr != nil {
		log.Println(err)
		return context.JSON(http.StatusInternalServerError, errors.New("Cannot delete user with email: "+email))
	}
	return context.NoContent(http.StatusNoContent)
}

// Login attempts to authenticate a given user
func Login(context *echo.Context) error {
	email, password := context.Form("email"), context.Form("password")
	if email == "" {
		return context.JSON(http.StatusBadRequest, errors.New("Missing email parameter in POST body"))
	}
	if password == "" {
		return context.JSON(http.StatusBadRequest, errors.New("Missing password parameter in POST body"))
	}

	user, err := LoadByEmail(email)

	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusUnauthorized, errors.New("Wrong email and password combination"))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusUnauthorized, errors.New("Wrong email and password combination"))
	}
	user.AuthToken = auth.GetEncodedToken()
	user.Hash = "" // dont return the hash, for security concerns

	userSession, err := session.NewSession(user.ID, user.AuthToken, user.IsAdmin)
	if err != nil {
		context.JSON(http.StatusInternalServerError, err)
	}
	session.Save(userSession)

	cookieExpiration := time.Now().Add(30 * 24 * time.Hour)
	cookie := http.Cookie{Name: "AuthToken", Value: user.AuthToken, Expires: cookieExpiration}
	context.Response().Header().Set("Set-Cookie", cookie.String())

	return context.JSON(http.StatusOK, user)
}

// ChangePassword attempts to authenticate a given user
func ChangePassword(context *echo.Context) error {
	email, oldPassword, newPassword := context.Form("email"), context.Form("oldPassword"), context.Form("newPassword")
	if email == "" {
		log.Println("Missing email")
		return context.JSON(http.StatusBadRequest, errors.New("Missing email parameter in POST body"))
	}
	if oldPassword == "" {
		log.Println("Missing oldPassword")
		return context.JSON(http.StatusBadRequest, errors.New("Missing oldPassword parameter in POST body"))
	}

	user, err := LoadByEmail(email)

	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusUnauthorized, errors.New("Wrong email and password combination"))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(oldPassword))
	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusUnauthorized, errors.New("Wrong email and password combination"))
	}
	// TODO factorise the above lines
	hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if bcryptErr != nil {
		log.Println(bcryptErr)
		return context.JSON(http.StatusInternalServerError, errors.New("Server error. The password was not changed."))
	}

	err = SavePassword(user, string(hash))
	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusInternalServerError, errors.New("Server error. The password was not changed."))
	}

	return context.JSON(http.StatusOK, nil)
}
