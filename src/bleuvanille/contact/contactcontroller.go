package contact

import (
	"bleuvanille/config"
	"fmt"
	"log"
	"net/http"
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
	CreatedAt string `json:"created_at"`
}

// LandingPage displays the landing page for getting new contacts
func LandingPage(context *echo.Context) error {

	_, err := context.Request().Cookie("visitorId")
	if err != nil {
		// first visit, we add a cookie
		fmt.Println("first visit")
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
	var contacts Contacts
	var err error
	switch sortParam {
	case "newer":
		contacts, err = LoadAll("created_at", "DESC")
	case "older":
		contacts, err = LoadAll("created_at", "ASC")
	case "email":
		contacts, err = LoadAll("email", "ASC")
	default:
		contacts, err = LoadAll("created_at", "DESC")
	}

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Contact list retrieval error"})
	}
	results := make([]formattedContact, len(contacts))
	for i := range contacts {
		formattedDate := monday.Format(contacts[i].CreatedAt, "Mon _2 Jan 2006 15:04", monday.LocaleFrFR)
		results[i] = formattedContact{contacts[i].Email, formattedDate}
		i++
	}
	return context.JSON(http.StatusOK, results)
}

// Create creates a new contact
func Create(context *echo.Context) error {

	email := context.Form("email")
	if email == "" {
		log.Println("Contact create email is null")
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing email parameter in POST body"})
	}
	timeSpent := context.Request().PostFormValue("timeSpent")
	// timeSpent := context.Form("timeSpent")
	if timeSpent != "" {
		timeSpentInt, err := strconv.Atoi(timeSpent)
		if err == nil {
			fmt.Printf("\ntimeSpent %v\n", timeSpentInt)
		} else {
			fmt.Println(err.Error())
		}
	}
	// TODO get time spent on the page
	// TODO check email is valid
	// TODO save these values
	// userAgent := context.Request().Header.Get("User-Agent")
	// referer := context.Request().Header.Get("Referer")

	contact, err := New(email)
	// contact, err := New(email, userAgent, referer)
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

// Remove deletes the contact for a given email
func Remove(context *echo.Context) error {
	email := context.Query("email")
	if email == "" {
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing email parameter in GET request"})
	}
	notFound, err := Delete(email)

	if err != nil {
		if notFound {
			return context.JSON(http.StatusNotFound, err)
		}
		log.Printf("Cannot delete contact with email %s, error: %s", email, err)
		return context.JSON(http.StatusInternalServerError, fmt.Errorf("Cannot delete contact with email %s", email))
	}
	return context.NoContent(http.StatusNoContent)
}
