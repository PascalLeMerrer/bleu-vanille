package contact

import (
	"bleuvanille/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/goodsign/monday"

	"github.com/labstack/echo"
	"github.com/twinj/uuid"
)

type errorMessage struct {
	Message string `json:"error"`
}

type formattedContact struct {
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UserAgent string `json:"userAgent"`
	Referer   string `json:"referer"`
	TimeSpent int    `json:"timeSpent"`
}

// LandingPage displays the landing page for getting new contacts
func LandingPage(context *echo.Context) error {

	_, err := context.Request().Cookie("visitorId")
	if err != nil {
		// first visit, we add a cookie
		expire := time.Now().AddDate(0, 0, 7) //7 days from now
		cookie := http.Cookie{
			Name:    "visitorId",
			Value:   uuid.NewV4().String(),
			Path:    "/",
			Domain:  config.HostName,
			Expires: expire,
		}
		http.SetCookie(context.Response().Writer(), &cookie)
		// TODO: increment unique visitor counter in database
	}
	return context.Render(http.StatusOK, "index", nil)
}

// GetAll writes the list of all contacts
func GetAll(context *echo.Context) error {
	sortParam := context.Query("sort")
	offsetParam, offsetErr := strconv.Atoi(context.Query("offset"))
	if offsetErr != nil {
		offsetParam = 0
	}
	limitParam, limitErr := strconv.Atoi(context.Query("limit"))
	if limitErr != nil {
		limitParam = 0
	}
	var contacts Contacts
	var totalCount int
	var err error
	switch sortParam {
	case "newer":
		contacts, totalCount, err = FindAll("created_at", "DESC", offsetParam, limitParam)
	case "older":
		contacts, totalCount, err = FindAll("created_at", "ASC", offsetParam, limitParam)
	case "emailAsc":
		contacts, totalCount, err = FindAll("email", "ASC", offsetParam, limitParam)
	case "emailDesc":
		contacts, totalCount, err = FindAll("email", "DESC", offsetParam, limitParam)
	default:
		contacts, totalCount, err = FindAll("created_at", "DESC", offsetParam, limitParam)
	}

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Contact list retrieval error"})
	}
	formattedContacts := make([]formattedContact, len(contacts))
	for i := range contacts {
		formattedDate := monday.Format(contacts[i].CreatedAt, "Mon _2 Jan 2006 15:04", monday.LocaleFrFR)
		formattedContacts[i] = formattedContact{contacts[i].Email, formattedDate, contacts[i].UserAgent, contacts[i].Referer, contacts[i].TimeSpent}
		i++
	}
	contentType := context.Request().Header.Get("Accept")
	if contentType != "" && len(contentType) >= len(echo.ApplicationJSON) && contentType[:len(echo.ApplicationJSON)] == echo.ApplicationJSON {
		context.Response().Header().Set("X-TOTAL-COUNT", strconv.Itoa(totalCount))
		return context.JSON(http.StatusOK, formattedContacts)
	}
	filepath, filename, err := createCsvFile(formattedContacts)
	if err != nil {
		fmt.Printf("Cannot create contact list file: %v", err)
		return context.JSON(http.StatusInternalServerError, fmt.Errorf("Cannot open file: %v", err))
	}
	// TODO: How to cleanup the temp dir?
	return context.File(filepath, filename, true)
}

// Create a CSV file containing the list of contacts
// returns the absolute file name (including the path) and the filename
func createCsvFile(formattedContacts []formattedContact) (string, string, error) {
	csvString := "Email, Date d'inscription, User Agent, Provenance, Temps passé sur landing page"
	for j := range formattedContacts {
		csvString += fmt.Sprintf("\"%s\", \"%s\", \"%s\", \"%s\", \"%d\"\n", formattedContacts[j].Email, formattedContacts[j].CreatedAt, formattedContacts[j].UserAgent, formattedContacts[j].Referer, formattedContacts[j].TimeSpent)
	}
	now := monday.Format(time.Now(), "2006-01-02-15h04", monday.LocaleFrFR)
	filename := "contactlist-" + now + ".csv"
	filepath := path.Join(os.TempDir(), filename)

	fileHandler, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return "", "", fmt.Errorf("Cannot open file: %v", err)
	}

	defer fileHandler.Close()

	_, err = fileHandler.Write([]byte(csvString))
	if err != nil {
		return "", "", fmt.Errorf("Cannot write file: %v", err)
	}
	return filepath, filename, nil
}

// Create creates a new contact
func Create(context *echo.Context) error {

	email := context.Form("email")
	if email == "" {
		log.Println("Contact create email is null")
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing email parameter in POST body"})
	}
	timeSpent := context.Request().PostFormValue("timeSpent")
	var timeSpentInt int
	// timeSpent := context.Form("timeSpent")
	if timeSpent != "" {
		var err error
		timeSpentInt, err = strconv.Atoi(timeSpent)
		if err != nil {
			log.Println(err.Error())
			timeSpentInt = -1
		}
	}
	// TODO check email is valid
	userAgent := context.Request().Header.Get("User-Agent")
	referer := context.Request().Header.Get("Referer")
	contact, err := New(email, userAgent, referer, timeSpentInt)
	if err != nil {
		log.Printf("Contact create error %v\n", err)
		return context.JSON(http.StatusInternalServerError, errorMessage{"Contact creation error"})
	}

	err = Save(&contact)
	if err != nil {
		if err.Error() == "cannot create document, unique constraint violated" {
			return context.JSON(http.StatusConflict, errorMessage{"Contact is already registered"})
		}
		log.Printf("Error: cannot save contact with email: %v\n", err)
		return context.JSON(http.StatusInternalServerError, errorMessage{"Contact creation error"})
	}

	return context.JSON(http.StatusCreated, contact)
}

// Delete deletes the contact for a given email
func Delete(context *echo.Context) error {
	email := context.Query("email")
	if email == "" {
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing email parameter in GET request"})
	}
	err := Remove(email)

	if err != nil {
		log.Printf("Cannot delete contact with email %s, error: %s", email, err)
		return context.JSON(http.StatusInternalServerError, fmt.Errorf("Cannot delete contact with email %s", email))
	}
	return context.NoContent(http.StatusNoContent)
}
