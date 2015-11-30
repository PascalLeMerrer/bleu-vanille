package eatable

import (
	"bleuvanille/log"

	//	"fmt"
	"net/http"
	//	"os"
	//	"path"
	//	"strconv"
	"time"
	//	"errors"
	"encoding/json"
	"io/ioutil"

	//	"github.com/goodsign/monday"

	"github.com/labstack/echo"
	//	"github.com/twinj/uuid"
)

type errorMessage struct {
	Message string `json:"error"`
}

var validEatableType map[string]string

func init() {
	auxmap := map[string]string{
		"ingredient":        "ingredient",
		"ingredientrecette": "ingredientrecette"}

	validEatableType = auxmap
}

//// LandingPage dismake(map[string]int)plays the landing page for getting new contacts
//func LandingPage(context *echo.Context) error {
//
//	_, err := context.Request().Cookie("visitorId")
//	if err != nil {
//		// first visit, we add a cookie
//		expire := time.Now().AddDate(0, 0, 7) //7 days from now
//		cookie := http.Cookie{
//			Name:    "visitorId",
//			Value:   uuid.NewV4().String(),
//			Path:    "/",
//			Domain:  config.HostName,
//			Expires: expire,
//		}
//		http.SetCookie(context.Response().Writer(), &cookie)
//		// TODO: increment unique visitor counter in database
//	}
//	return context.Render(http.StatusOK, "index", nil)
//}

//Return the object stored in database
func Get(context *echo.Context) error {
	id := context.Param("id")

	result, error := LoadById(id)

	//Verify if the result is correctly retreive from database
	if error != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Unvalid ID : " + id})
	}

	//The temporary struct is used to remove the fields _id, _rev and _key and add id field
	return context.JSON(http.StatusOK, struct {
		*Eatable
		Id      string `json:"id,omitempty"`
		OmitId  omit   `json:"_id,omitempty"`
		OmitRev omit   `json:"_rev,omitempty"`
		OmitKey omit   `json:"_key,omitempty"`
	}{Eatable: result, Id: result.Id})
}

//// GetAll writes the list of all contacts
//func GetAll(context *echo.Context) error {
//	sortParam := context.Query("sort")
//	var contacts Contacts
//	var err error
//	switch sortParam {
//	case "newer":
//		contacts, err = LoadAll("created_at", "DESC")
//	case "older":
//		contacts, err = LoadAll("created_at", "ASC")
//	case "email":
//		contacts, err = LoadAll("email", "ASC")
//	default:
//		contacts, err = LoadAll("created_at", "DESC")
//	}
//
//	if err != nil {
//		return context.JSON(http.StatusInternalServerError, errorMessage{"Contact list retrieval error"})
//	}
//	formattedContacts := make([]formattedContact, len(contacts))
//	for i := range contacts {
//		formattedDate := monday.Format(contacts[i].CreatedAt, "Mon _2 Jan 2006 15:04", monday.LocaleFrFR)
//		formattedContacts[i] = formattedContact{contacts[i].Email, formattedDate, contacts[i].UserAgent, contacts[i].Referer, contacts[i].TimeSpent}
//		i++
//	}
//	contentType := context.Request().Header.Get(echo.ContentType)
//	if contentType != "" && len(contentType) >= len(echo.ApplicationJSON) && contentType[:len(echo.ApplicationJSON)] == echo.ApplicationJSON {
//		return context.JSON(http.StatusOK, formattedContacts)
//	}
//	filepath, filename, err := createCsvFile(formattedContacts)
//	if err != nil {
//		fmt.Printf("Cannot create contact list file: %v", err)
//		return context.JSON(http.StatusInternalServerError, fmt.Errorf("Cannot open file: %v", err))
//	}
//	// TODO: How to cleanup the temp dir?
//	return context.File(filepath, filename, true)
//}

type omit *struct{}

// Create creates a new contact
func Create(context *echo.Context) error {

	bodyio := context.Request().Body

	bodybytes, err := ioutil.ReadAll(bodyio)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Error while reading the eatable content"})
	}

	if len(bodybytes) == 0 {
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing eatable content"})
	}

	eatable := NewEmpty()

	//To ignore the status, created_at and nutrient information, we create fake fields that will not be stored in the eatable object itself
	err = json.Unmarshal(bodybytes,
		&struct {
			*Eatable
			Status    string    `json:"status"`
			CreatedAt time.Time `json:"created_at"`
			Nutrient  *Nutrient `json:"nutrient"`
		}{Eatable: eatable})

	if err != nil {
		return context.JSON(http.StatusCreated, errorMessage{"Incorrect eatable content : " + err.Error()})
	}

	//Check the type of eatable
	if eatable.Type != "" {
		existingType := validEatableType[eatable.Type]
		if existingType == "" {
			return context.JSON(http.StatusBadRequest, errorMessage{"Unknow type : " + eatable.Type})
		}
	}

	//Save the new eatable object
	err = Save(eatable)

	if err != nil {
		log.Printf("Error: cannot save eatable : %v\n", err)
		return context.JSON(http.StatusInternalServerError, errorMessage{"Eatable creation error"})
	}

	return context.JSON(http.StatusCreated, eatable)
}

// Update updates an existing contact
func Update(context *echo.Context) error {
	id := context.Param("id")

	//Read the body that contains the new json for our eatable
	bodyio := context.Request().Body

	bodybytes, err := ioutil.ReadAll(bodyio)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Error while reading the eatable content"})
	}

	//body should not be empty
	if len(bodybytes) == 0 {
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing eatable content"})
	}

	//Read the eatable from database
	eatable, errload := LoadById(id)

	if errload != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Incorrect eatable id : " + err.Error()})
	}

	//Populate the existing eatable with the new data, ignoring, status, created time and nutrient
	err = json.Unmarshal(bodybytes, &struct {
		*Eatable
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		Nutrient  *Nutrient `json:"nutrient"`
	}{Eatable: eatable})

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Incorrect eatable content : " + err.Error()})
	}

	//Check the type of eatable
	if eatable.Type != "" {
		existingType := validEatableType[eatable.Type]
		if existingType == "" {
			return context.JSON(http.StatusBadRequest, errorMessage{"Unknow type : " + eatable.Type})
		}
	}

	//Save the new eatable object
	err = Save(eatable)

	if err != nil {
		log.Printf("Error: cannot update eatable : %v\n", err)
		return context.JSON(http.StatusInternalServerError, errorMessage{"Update eatable error"})
	}

	return context.JSON(http.StatusOK, eatable)
}

// Update updates an existing contact
func SetNutrient(context *echo.Context) error {
	id := context.Param("id")

	bodyio := context.Request().Body

	bodybytes, err := ioutil.ReadAll(bodyio)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Error while reading the eatable content"})
	}

	if len(bodybytes) == 0 {
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing eatable content"})
	}

	eatable, errload := LoadById(id)

	if errload != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Incorrect eatable id : " + err.Error()})
	}

	var nutrient Nutrient
	err = json.Unmarshal(bodybytes, &nutrient)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Incorrect nutrient content : " + err.Error()})
	}

	eatable.Nutrient = &nutrient

	err = Save(eatable)

	if err != nil {
		log.Printf("Error: cannot set nutrient : %v\n", err)
		return context.JSON(http.StatusInternalServerError, errorMessage{"Set Nutrient Eatable error"})
	}

	return context.JSON(http.StatusOK, eatable)
}

// Update updates an existing contact
func SetParent(context *echo.Context) error {
	id := context.Param("id")
	idParent := context.Param("newParentId")

	err := SaveParent(id, idParent)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Impossible to set parent : " + err.Error()})
	}

	return context.JSON(http.StatusOK, "ok")
}
