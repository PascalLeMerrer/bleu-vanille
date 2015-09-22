package user

import (
	"bleuvanille/auth"
	"bleuvanille/config"
	"bleuvanille/mail"
	"bleuvanille/session"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
)

// CreateDefault creates a default admin account if it does not exist
func CreateDefault() {
	existingAdmin, err := LoadByEmail(config.AdminEmail)
	if err != nil {
		log.Println(err)
	}

	if !existingAdmin.IsAdmin {
		admin, err := New("admin@bleuvanille.com", "Admin", "Admin", "xeCuf8CHapreNe=")
		if err != nil {
			log.Fatal(err)
		}
		admin.IsAdmin = true
		err = Save(admin)
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				log.Fatal("User admin@bleuvanille.com but has no admin rights. WTF?!")
			}
		}
		log.Println("Admin account created with default password. You should change it.")
		return
	}
	log.Println("Admin account found.")
}

// Create creates a new user
var Create = emailAndPasswordRequired(
	func(context *echo.Context) error {

		email := context.Form("email")
		password := context.Form("password")
		firstname := context.Form("firstname")
		lastname := context.Form("lastname")

		if firstname == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Missing firstname parameter in POST body"))
		}
		if lastname == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Missing lastname parameter in POST body"))
		}

		user, err := New(email, firstname, lastname, password)
		if err != nil {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, errors.New("User creation error"))
		}

		err = Save(user)
		if err != nil {
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
	})

// Get returns the entry for a given email if any
var Get = emailRequired(
	func(context *echo.Context) error {

		email := context.Query("email")
		var user, err = LoadByEmail(email)
		if err != nil {
			log.Println(err)
			return context.JSON(http.StatusNotFound, nil)
		}
		return context.JSON(http.StatusOK, user)
	})

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

	userID := context.Get("session").(*session.Session).UserID
	if userID == "" {
		log.Println("Missing userID in session")
		return context.JSON(http.StatusUnauthorized, errors.New("Maybe your session expired. Try to disconnect then reconnect."))
	}
	password := context.Form("password")
	if password == "" {
		log.Println("Missing password parameter in DELETE request")
		return context.JSON(http.StatusBadRequest, errors.New("Missing password parameter in DELETE request"))
	}
	user, err := LoadByID(userID)

	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusUnauthorized, errors.New("Wrong email and password combination"))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		return context.JSON(http.StatusUnauthorized, errors.New("Wrong email and password combination"))
	}

	deleteErr := Delete(*user)
	if deleteErr != nil {
		log.Println(err)
		return context.JSON(http.StatusInternalServerError, errors.New("Cannot delete user with ID: "+userID))
	}
	return context.NoContent(http.StatusNoContent)
}

// Login attempts to authenticate a given user
var Login = emailAndPasswordRequired(
	func(context *echo.Context) error {
		password := context.Form("password")
		email := context.Form("email")
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
		authToken, sessionID := auth.GetSessionToken()
		user.Hash = "" // dont return the hash, for security concerns

		userSession, err := session.New(sessionID, user.ID, user.IsAdmin)
		if err != nil {
			context.JSON(http.StatusInternalServerError, err)
		}
		session.Save(userSession)

		context.Response().Header().Set(echo.Authorization, authToken)

		return context.JSON(http.StatusOK, user)
	})

// ChangePassword updates the password of the authenticated user
var ChangePassword = emailAndPasswordRequired(
	func(context *echo.Context) error {

		// TODO we should use info from session instead of an email parameter
		email, oldPassword, newPassword := context.Form("email"), context.Form("password"), context.Form("newPassword")

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
	})

// ResetPassword updates the password of the authenticated user without proving the old one
func ResetPassword(context *echo.Context) error {
	log.Println("ResetPassword")
	newPassword := context.Form("password")
	var email string
	var ok bool
	if email, ok = context.Get("email").(string); !ok {
		log.Println("Missing email in password reset request")
		return context.JSON(http.StatusUnauthorized, errors.New("Missing email in password reset request"))
	}
	user, err := LoadByEmail(email)
	if err != nil {
		log.Printf("Cannot find user during password reset: %v\n", err)
		return context.JSON(http.StatusUnauthorized, errors.New("Wrong email in password reset request"))
	}

	hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if bcryptErr != nil {
		log.Println(bcryptErr)
		return context.JSON(http.StatusInternalServerError, errors.New("Server error. The password was not changed."))
	}

	user.Hash = string(hash)
	err = Update(*user)
	if err != nil {
		log.Println(err)
		return context.JSON(http.StatusInternalServerError, errors.New("Server error. The password was not changed."))
	}

	return context.JSON(http.StatusOK, nil)
}

// SendResetLink send an email to reset the password of the current User
var SendResetLink = emailRequired(
	func(context *echo.Context) error {
		email := context.Form("email")
		user, err := LoadByEmail(email)
		if err != nil {
			log.Println(err)
			return context.JSON(http.StatusNotFound, errors.New("Cannot find user for email "+email))
		}
		user.ResetToken = auth.GetResetToken(email)
		err = Update(*user)
		if err != nil {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, errors.New("Cannot save password reset token for user "+email))
		}
		// TODO Use a template
		resetLink := fmt.Sprintf("To: %v\r\n"+
			"Subject: Password reset\r\n"+
			"\r\n<!DOCTYPE html><html><body><a href=\"http://%v:%v/users/resetform?token=%v\">Reset Password</a></body></html>",
			email,
			config.HostName,
			config.HostPort,
			user.ResetToken)
		log.Println(resetLink)
		err = mail.Send(email, []byte(resetLink))
		if err != nil {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, errors.New("Cannot send password reset email to user "+email))
		}

		return context.JSON(http.StatusOK, "Password reset email sent to "+email)
	})

// DisplayResetForm displays the reset password form
func DisplayResetForm(context *echo.Context) error {
	token := context.Query("token")
	if token == "" {
		return context.JSON(http.StatusUnauthorized, "Invalid URL for password reset")

	}
	data := struct {
		Token string
	}{
		Token: token,
	}
	return context.Render(http.StatusOK, "passwordreset", data)
}

// Profile Returns the profile of a given user
// If it's the current one, the returned profile is richer
func Profile(context *echo.Context) error {
	//TODO implement
	return context.JSON(http.StatusOK, "user profile")
}

// emailRequired verifies the presence of an "email" parameter in the body of the request
// this function is a kind of decorator used to wrap an echo handler function
// (http://talks.golang.org/2013/go4python.slide#37)
func emailRequired(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(context *echo.Context) error {
		email := context.Form("email")
		if email == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Missing email parameter in POST body"))
		}
		return handler(context)
	}
}

// emailAndPasswordRequired verifies the presence of an "email" parameter
// and a "password" parameter in the body of the request
// this function is a kind of decorator used to wrap an echo handler function
// (http://talks.golang.org/2013/go4python.slide#37)
func emailAndPasswordRequired(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(context *echo.Context) error {
		email := context.Form("email")
		if email == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Missing email parameter in POST body"))
		}
		password := context.Form("password")
		if password == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Missing password parameter in POST body"))
		}
		return handler(context)
	}
}
