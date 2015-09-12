package contact

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type errorMessage struct {
	Message string `json:"error"`
}

// GetAll writes the list of all contacts
func GetAll(context *echo.Context) error {
	contacts, err := LoadAll()
	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Contact list retrieval error"})
	}
	return context.JSON(http.StatusOK, contacts)
}

// Create creates a new contact
func Create(context *echo.Context) error {

	email := context.Form("email")
	if email == "" {
		log.Println("Contact create email is null")
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing email parameter in POST body"})
	}
	//TODO check email is valid
	contact, err := NewContact(email)
	if err != nil {
		log.Printf("Contact create error %v\n", err)
		return context.JSON(http.StatusInternalServerError, errorMessage{"Contact creation error"})
	}

	_, err = Save(contact)
	if err != nil {
		log.Printf("Contact save error %v\n", err)
		return context.JSON(http.StatusInternalServerError, errorMessage{"Contact creation error"})
	}

	return context.JSON(http.StatusCreated, contact)
}

// Get returns the entry for a given email if any
func Get(context *echo.Context) error {

	email := context.Query("email")
	if email == "" {
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing email parameter in GET request"})
	}

	var contact, err = LoadByEmail(email)
	if err != nil {
		return context.NoContent(http.StatusNotFound)
	}
	return context.JSON(http.StatusOK, contact)
}

// Remove deletes the contact for a given email
func Remove(context *echo.Context) error {
	email := context.Query("email")
	if email == "" {
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing email parameter in GET request"})
	}
	err := Delete(email)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Cannot delete contact with email: " + email})
	}
	return context.NoContent(http.StatusNoContent)
}
