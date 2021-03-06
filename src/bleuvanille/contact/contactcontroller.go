package contact

import (
	"bleuvanille/config"
	"bleuvanille/statistics"
	"errors"
	"fmt"
	"github.com/goodsign/monday"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/goware/emailx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/twinj/uuid"
)

type formattedContact struct {
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UserAgent string `json:"userAgent"`
	Referer   string `json:"referer"`
	TimeSpent int    `json:"timeSpent"`
}

// LandingPage displays the landing page for getting new contacts
func LandingPage() echo.HandlerFunc {
	return func(context echo.Context) error {

		_, err := context.Request().(*standard.Request).Request.Cookie("visitorId")
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
			context.Response().Header().Set("Set-Cookie", cookie.String())

			// increment unique visitor counter in database
			statistics.IncrementCounter("landing_page_visitor_count")
		}
		return context.Render(http.StatusOK, "index", nil)
	}
}

// GetAll writes the list of all contacts
func GetAll() echo.HandlerFunc {
	return func(context echo.Context) error {
		sortParam := context.QueryParam("sort")
		offsetParam, offsetErr := strconv.Atoi(context.QueryParam("offset"))
		if offsetErr != nil {
			offsetParam = 0
		}
		limitParam, limitErr := strconv.Atoi(context.QueryParam("limit"))
		if limitErr != nil {
			limitParam = 0
		}
		var contacts []Contact
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
			return context.JSON(http.StatusInternalServerError, errors.New("Contact list retrieval error"))
		}
		formattedContacts := make([]formattedContact, len(contacts))
		for i := range contacts {
			formattedDate := monday.Format(contacts[i].CreatedAt, "Mon _2 Jan 2006 15:04", monday.LocaleFrFR)
			formattedContacts[i] = formattedContact{contacts[i].Email, formattedDate, contacts[i].UserAgent, contacts[i].Referer, contacts[i].TimeSpent}
			i++
		}
		contentType := context.Request().Header().Get("Accept")
		if contentType != "" && len(contentType) >= len(echo.MIMEApplicationJSON) && contentType[:len(echo.MIMEApplicationJSON)] == echo.MIMEApplicationJSON {
			context.Response().Header().Set("X-TOTAL-COUNT", strconv.Itoa(totalCount))
			return context.JSON(http.StatusOK, formattedContacts)
		}
		filepath, err := createCsvFile(formattedContacts)
		if err != nil {
			fmt.Printf("Cannot create contact list file: %v", err)
			return context.JSON(http.StatusInternalServerError, fmt.Errorf("Cannot open file: %v", err))
		}
		return context.File(filepath)
	}
}

// Create a CSV file containing the list of contacts
// returns the absolute file name
func createCsvFile(formattedContacts []formattedContact) (string, error) {
	csvString := "Email, Date d'inscription, User Agent, Provenance, Temps passé sur landing page"
	for j := range formattedContacts {
		csvString += fmt.Sprintf("\"%s\", \"%s\", \"%s\", \"%s\", \"%d\"\n", formattedContacts[j].Email, formattedContacts[j].CreatedAt, formattedContacts[j].UserAgent, formattedContacts[j].Referer, formattedContacts[j].TimeSpent)
	}
	now := monday.Format(time.Now(), "2006-01-02-15h04", monday.LocaleFrFR)
	filename := "contactlist-" + now + ".csv"
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

// Create creates a new contact
func Create() echo.HandlerFunc {
	return func(context echo.Context) error {
		email := context.FormValue("email")
		if email == "" {
			log.Println("Contact creation email is null")
			return context.JSON(http.StatusBadRequest, errors.New("Missing email parameter in POST body"))
		}

		err := emailx.Validate(email)
		if err != nil {
			return context.JSON(http.StatusBadRequest, errors.New("Invalid email parameter in POST body"))
		}

		timeSpent := context.Request().FormValue("timeSpent")
		var timeSpentInt int
		if timeSpent != "" {
			var err error
			timeSpentInt, err = strconv.Atoi(timeSpent)
			if err != nil {
				log.Println(err.Error())
				timeSpentInt = -1
			}
		}
		userAgent := context.Request().UserAgent()
		referer := context.Request().Header().Get("Referer")
		contact, err := New(email, userAgent, referer, timeSpentInt)
		if err != nil {
			log.Printf("Contact create error %v\n", err)
			return context.JSON(http.StatusInternalServerError, errors.New("Contact creation error"))
		}

		err = Save(&contact)
		if err != nil {
			if err.Error() == "cannot create document, unique constraint violated" {
				return context.JSON(http.StatusConflict, errors.New("Contact is already registered"))
			}
			log.Printf("Error: cannot save contact with email: %v\n", err)
			return context.JSON(http.StatusInternalServerError, errors.New("Contact creation error"))
		}
		return context.JSON(http.StatusCreated, contact)
	}
}

// Delete deletes the contact for a given email
func Delete() echo.HandlerFunc {
	return func(context echo.Context) error {
		email := context.QueryParam("email")
		if email == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Missing email parameter in GET request"))
		}
		err := Remove(email)

		if err != nil {
			log.Printf("Cannot delete contact with email %s, error: %s", email, err)
			return context.JSON(http.StatusInternalServerError, fmt.Errorf("Cannot delete contact with email %s", email))
		}
		return context.NoContent(http.StatusNoContent)
	}
}
