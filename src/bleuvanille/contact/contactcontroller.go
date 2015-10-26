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
	contact, err := New(email)
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
	err := Delete(email)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Cannot delete contact with email: " + email})
	}
	return context.NoContent(http.StatusNoContent)
}
