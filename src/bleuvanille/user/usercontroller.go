package user

import (
	"bleuvanille/auth"
	"bleuvanille/config"
	"bleuvanille/mail"
	"bleuvanille/session"
	"bytes"
	"errors"
	"fmt"
	"github.com/goodsign/monday"
	"github.com/goware/emailx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo"
)

// data used to fill the password reset email template
type passwordResetData struct {
	From     string
	To       string
	Host     string
	Port     int
	Token    string
	Now      string
	Boundary string
}

type formattedUser struct {
	ID        string `json:"id"`
	Email     string `json:"email,omitempty"`
	CreatedAt string `json:"createdAt"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	IsAdmin   bool   `json:"isAdmin"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// CreateDefault creates a default admin account if it does not exist
func CreateDefault() {
	existingAdmin, err := FindByEmail(config.AdminEmail)
	if err != nil {
		log.Println(err)
	}
	if existingAdmin != nil && !existingAdmin.IsAdmin {
		log.Fatalf("FATAL: User with email %v exists but has no admin rights.", config.AdminEmail)
	}

	if existingAdmin == nil {
		admin, err := New(config.AdminEmail, "Admin", "Admin", config.AdminPassword)
		if err != nil {
			log.Fatal(err)
		}
		admin.IsAdmin = true
		err = Save(&admin)
		if err != nil {
			log.Fatalf("Cannot create admin user with email %v. Error: %v", config.AdminEmail, err.Error())
		}
		return
	}
	if config.ServerDebug() {
		log.Println("Admin account found.")
	}

	//Create the test user if we are not in production
	if !config.ProductionMode {
		existingTestUser, err := FindByEmail(config.TestUserEmail)

		if err != nil {
			log.Println(err)
		}

		if existingTestUser == nil {
			testuser, err := New(config.TestUserEmail, "TestUser", "TestUser", "xaFqJDeJldIEcdfZS")
			if err != nil {
				log.Fatal(err)
			}
			err = Save(&testuser)
			if err != nil {
				log.Fatalf("Cannot create user with email %v. Error: %v", config.TestUserEmail, err.Error())
			}
			return
		}
		log.Println("Test account found.")
	}
}

// Create creates a new user
func Create() echo.HandlerFunc {
	return emailAndPasswordRequired(
		func(context echo.Context) error {

			email := context.FormValue("email")
			password := context.FormValue("password")
			firstname := context.FormValue("firstname")
			lastname := context.FormValue("lastname")

			err := emailx.Validate(email)
			if err != nil {
				return context.JSON(http.StatusBadRequest, errors.New("Invalid email parameter in POST body"))
			}

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

			err = Save(&user)
			if err != nil {
				if err.Error() == "cannot create document, unique constraint violated" {
					return context.JSON(http.StatusConflict, errors.New("User already exists"))
				}
				return context.JSON(http.StatusInternalServerError, errors.New("User account creation failed - cannot be saved"))
			}
			user.Hash = "" // dont return the hash, for security concerns

			return context.JSON(http.StatusCreated, formatUser(&user))
		})
}

// converts a User object to a FormattedUSer object
// which is its public counterpart
func formatUser(user *User) formattedUser {
	return formattedUser{user.ID, user.Email, formatDate(user.CreatedAt), user.Firstname, user.Lastname, user.IsAdmin}
}

// Get returns the entry for a given email if any
func Get() echo.HandlerFunc {
	return emailRequired(
		func(context echo.Context) error {

			email := context.QueryParam("email")
			var user, err = FindByEmail(email)
			if err != nil || user == nil {
				log.Println(err)
				return context.JSON(http.StatusNotFound, nil)
			}
			return context.JSON(http.StatusOK, formatUser(user))
		})
}

// GetAll returns the list of all users
func GetAll() echo.HandlerFunc {
	return func(context echo.Context) error {
		offsetParam, offsetErr := strconv.Atoi(context.QueryParam("offset"))
		if offsetErr != nil {
			offsetParam = 0
		}
		limitParam, limitErr := strconv.Atoi(context.QueryParam("limit"))
		if limitErr != nil {
			limitParam = 0
		}
		var users []User
		var totalCount int
		var err error
		criteria, order := getSortingParams(context)

		users, totalCount, err = FindAll(criteria, order, offsetParam, limitParam)

		if err != nil {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, errors.New("User list retrieval error"))
		}
		formattedUsers := make([]formattedUser, len(users))
		for i := range users {
			formattedUsers[i] = formatUser(&users[i])
			i++
		}
		contentType := context.Request().Header().Get("Accept")
		if contentType != "" && len(contentType) >= len(echo.MIMEApplicationJSON) && contentType[:len(echo.MIMEApplicationJSON)] == echo.MIMEApplicationJSON {
			context.Response().Header().Set("X-TOTAL-COUNT", strconv.Itoa(totalCount))
			return context.JSON(http.StatusOK, formattedUsers)
		}
		filepath, err := createCsvFile(formattedUsers)
		if err != nil {
			fmt.Printf("Cannot create User list file: %v", err)
			return context.JSON(http.StatusInternalServerError, fmt.Errorf("Cannot open file: %v", err))
		}
		return context.File(filepath)
	}
}

// extracts from the sortParam parameter a criteria value and a sorting order
// @returns the nma of the property on which user list must be sorted, and the sort order (ASC or DESC)
func getSortingParams(context echo.Context) (string, string) {

	sortParam := context.QueryParam("sort")

	var criteria string
	var order string

	switch sortParam {
	case "newer":
		criteria = "createdAt"
		order = "DESC"
	case "older":
		criteria = "createdAt"
		order = "ASC"
	case "emailAsc":
		criteria = "email"
		order = "ASC"
	case "emailDesc":
		criteria = "email"
		order = "DESC"
	case "nameAsc":
		criteria = "lastname"
		order = "ASC"
	case "nameDesc":
		criteria = "lastname"
		order = "DESC"
	default:
		criteria = "createdAt"
		order = "DESC"
	}
	return criteria, order
}

// formats a date according to FR locale
func formatDate(date time.Time) string {
	return monday.Format(date, "Mon _2 Jan 2006 15:04", monday.LocaleFrFR)
}

// Create a CSV file containing the list of users
// returns the absolute file name
func createCsvFile(formattedUsers []formattedUser) (string, error) {
	csvString := "Email, Date d'inscription, Prénom, Nom, Admin"
	for j := range formattedUsers {
		csvString += fmt.Sprintf("\"%s\", \"%s\", \"%s\", \"%s\", \"%t\"\n", formattedUsers[j].Email, formattedUsers[j].CreatedAt, formattedUsers[j].Firstname, formattedUsers[j].Lastname, formattedUsers[j].IsAdmin)
	}
	now := monday.Format(time.Now(), "2006-01-02-15h04", monday.LocaleFrFR)
	filename := "Userlist-" + now + ".csv"
	filepath := path.Join(os.TempDir(), filename)

	fileHandler, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return "", fmt.Errorf("Cannot open file: %v", err)
	}

	defer fileHandler.Close()

	_, err = fileHandler.Write([]byte(csvString))
	if err != nil {
		return "", fmt.Errorf("Cannot write file: %v", err)
	}
	return filepath, nil
}

// Patch modifies the user account for a given ID
// This is an admin feature, not supposed to be used by normal users
func Patch() echo.HandlerFunc {
	return func(context echo.Context) error {
		userID := context.Param("userID")
		session := context.Get("session").(*session.Session)
		if session == nil || (session.UserID != userID && !session.IsAdmin) {
			log.Printf("ERROR: unauthorized attempt to modify account %s by  user with session %+v", userID, session)
			return context.JSON(http.StatusUnauthorized, "")
		}

		user, err := FindByID(userID)

		if err != nil || user == nil {
			return context.JSON(http.StatusInternalServerError, errors.New("Cannot load user with ID %s"))
		}

		previousIsAdmin := user.IsAdmin

		err = context.Bind(&user)
		if err != nil {
			log.Printf("Cannot bind user %v", err)
			return context.JSON(http.StatusBadRequest, errors.New("Cannot decode request body"))
		}
		if user.IsAdmin != previousIsAdmin && !session.IsAdmin {
			log.Printf("ERROR: unauthorized attempt to give admin rights to account %s by user with session %+v", user.Email, session)
			return context.JSON(http.StatusUnauthorized, "")
		}

		saveErr := Save(user)
		if saveErr != nil {
			log.Printf("Cannot update user %v", user.ID)
			return context.JSON(http.StatusInternalServerError, errors.New("Cannot update user "+user.ID))
		}
		user.Hash = "" // never leak the hash
		return context.JSON(http.StatusOK, user)
	}
}

// RemoveByAdmin removes the user account for a given ID
// This is an admin feature, not supposed to be used by normal users
func RemoveByAdmin() echo.HandlerFunc {
	return func(context echo.Context) error {
		user, err := FindByID(context.Param("userID"))
		if err != nil || user == nil {
			return context.JSON(http.StatusInternalServerError, errors.New("Cannot load user with ID %s"))

		}
		return _delete(context, user)
	}
}

// Delete removes the user account for a given email
func Delete() echo.HandlerFunc {
	return func(context echo.Context) error {
		userID := context.Get("session").(*session.Session).UserID
		if userID == "" {
			log.Println("Missing userID in session")
			return context.JSON(http.StatusUnauthorized, errors.New("Maybe your session expired. Try to disconnect then reconnect."))
		}
		password := context.FormValue("password")
		if password == "" {
			log.Println("Missing password parameter in DELETE request")
			return context.JSON(http.StatusBadRequest, errors.New("Missing password parameter in DELETE request"))
		}
		user, err := FindByID(userID)
		if err != nil || user == nil {
			log.Printf("Cannot load user with ID %s", userID)
			return context.JSON(http.StatusUnauthorized, errors.New("Wrong email and password combination"))
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
		if err != nil {
			return context.JSON(http.StatusUnauthorized, errors.New("Wrong email and password combination"))
		}
		return _delete(context, user)
	}
}

// Effective deletion of a user account
func _delete(context echo.Context, user *User) error {
	deleteErr := Remove(user.Key)
	if deleteErr != nil {
		log.Printf("Cannot delete user %v", deleteErr)
		return context.JSON(http.StatusInternalServerError, errors.New("Cannot delete user with ID: "+user.ID))
	}
	return context.NoContent(http.StatusNoContent)
}

// Login attempts to authenticate a given user
func Login() echo.HandlerFunc {
	return emailAndPasswordRequired(
		func(context echo.Context) error {
			password := context.FormValue("password")
			email := context.FormValue("email")
			user, err := FindByEmail(email)
			if err != nil || user == nil {
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

			userSession, err := session.New(sessionID, user.ID, user.Firstname, user.Lastname, user.Email, user.IsAdmin)
			if err != nil {
				context.JSON(http.StatusInternalServerError, err)
			}
			err = session.Save(&userSession)
			if err != nil {
				log.Printf("Cannot save session %+v: %s\n", userSession, err)
				context.JSON(http.StatusInternalServerError, err)
			}
			context.Response().Header().Set(echo.HeaderAuthorization, authToken)
			addCookie(context, authToken)
			return context.JSON(http.StatusOK, formatUser(user))
		})
}

// adds a session cookie
func addCookie(context echo.Context, authToken string) {
	expire := time.Now().AddDate(0, 1, 0) // 1 month
	cookie := &http.Cookie{
		Name:    "token",
		Expires: expire,
		Value:   auth.Bearer + " " + authToken,
		Path:    "/",
		// Domain must not be set for auth to work with chrome without domain name
		// http://stackoverflow.com/questions/5849013/setcookie-does-not-set-cookie-in-google-chrome
	}
	context.Response().Header().Set("Set-Cookie", cookie.String())
}

// ResetPassword updates the password of the authenticated user without providing the old one
func ResetPassword() echo.HandlerFunc {
	return func(context echo.Context) error {
		newPassword := context.FormValue("password")
		var email string
		var ok bool
		if email, ok = context.Get("email").(string); !ok {
			log.Println("Missing email in password reset request")
			return context.JSON(http.StatusUnauthorized, errors.New("Missing email in password reset request"))
		}
		user, err := FindByEmail(email)
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
		err = Save(user)
		if err != nil {
			log.Println(err)
			return context.JSON(http.StatusInternalServerError, errors.New("Server error. The password was not changed."))
		}

		return context.JSON(http.StatusOK, nil)
	}
}

// SendResetLink send an email to reset the password of the current User
func SendResetLink() echo.HandlerFunc {
	return emailRequired(
		func(context echo.Context) error {
			email := context.FormValue("email")
			user, err := FindByEmail(email)
			if err != nil || user == nil {
				log.Println(err)
				return context.JSON(http.StatusNotFound, errors.New("Cannot find user for email "+email))
			}
			user.ResetToken = auth.GetResetToken(email)
			err = Save(user)
			if err != nil {
				log.Println(err)
				return context.JSON(http.StatusInternalServerError, errors.New("Cannot save password reset token for user "+email))
			}

			data := passwordResetData{config.NoReplyAddress, email, config.HostName, config.HostPort, user.ResetToken, time.Now().String(), randomString(16)}
			emailTemplateTree := template.New("resetEmailTemplate")
			emailTemplateCollection := template.Must(emailTemplateTree.ParseFiles("src/bleuvanille/templates/PasswordReset.email"))

			var emailBody = new(bytes.Buffer)

			t := emailTemplateCollection.Templates()[0]

			err = t.Execute(emailBody, data)

			err = mail.Send(email, emailBody.Bytes())
			if err != nil {
				log.Println(err)
				return context.JSON(http.StatusInternalServerError, errors.New("Cannot send password reset email to user "+email))
			}

			return context.JSON(http.StatusOK, "Password reset email sent to "+email)
		})
}

// DisplayResetForm displays the reset password form
func DisplayResetForm() echo.HandlerFunc {
	return func(context echo.Context) error {
		token := context.QueryParam("token")
		email := context.QueryParam("email")
		if token == "" || email == "" {
			return context.JSON(http.StatusUnauthorized, "Invalid URL for password reset")
		}
		user, err := FindByEmail(email)
		if user == nil || err != nil {
			return context.JSON(http.StatusUnauthorized, "Invalid user for password reset")
		}

		data := struct {
			Token   string
			IsAdmin bool
		}{
			Token:   token,
			IsAdmin: user.IsAdmin,
		}
		return context.Render(http.StatusOK, "passwordreset", data)
	}
}

// Profile Returns the profile of a given user
// If it's the current one,or if if it's an admin that requires it,
// the returned profile is richer than if it a normal user that asks
func Profile() echo.HandlerFunc {
	return func(context echo.Context) error {
		user, err := FindByID(context.Param("userID"))
		if err != nil || user == nil {
			return context.JSON(http.StatusInternalServerError, errors.New("Cannot load user with ID %s"))
		}

		session := context.Get("session").(*session.Session)

		if session != nil && (session.IsAdmin || session.UserID == user.ID) {
			user.Hash = "" // don't leak user Hash, for security
			return context.JSON(http.StatusOK, user)
		}

		return context.JSON(http.StatusOK, formatUser(user))
	}
}

// emailRequired verifies the presence of an "email" parameter in the body of the request
// this function is a kind of decorator used to wrap an echo handler function
// (http://talks.golang.org/2013/go4python.slide#37)
func emailRequired(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		email := context.FormValue("email")
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
	return func(context echo.Context) error {
		email := context.FormValue("email")
		if email == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Missing email parameter in POST body"))
		}
		password := context.FormValue("password")
		if password == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Missing password parameter in POST body"))
		}
		return handler(context)
	}
}

// ChangePassword updates the password of the authenticated user
func ChangePassword() echo.HandlerFunc {
	return func(context echo.Context) error {
		session := context.Get("session").(*session.Session)
		if session == nil {
			log.Printf("ERROR: invalid session %+v", session)
			return context.JSON(http.StatusUnauthorized, "")
		}

		oldPassword, newPassword := context.FormValue("password"), context.FormValue("newPassword")
		if oldPassword == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Missing password parameter in POST body"))
		}

		email := session.Email
		user, err := FindByEmail(email)
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
	}
}

//TODO extract in a util package
// source: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randomString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
