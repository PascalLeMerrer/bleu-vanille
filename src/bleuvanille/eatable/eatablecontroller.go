package eatable

import (
	"bleuvanille/date"
	"bleuvanille/log"
	"encoding/json"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"time"
)

type errorMessage struct {
	Message string `json:"error"`
}

// Get returns the eatable object stored in database
func Get(context *echo.Context) error {
	key := context.Param("key")

	eatable, error := FindByKey(key)

	if error != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Invalid key: " + key})
	}

	if eatable == nil {
		return context.JSON(http.StatusNotFound, errorMessage{"No eatable found for key: " + key})
	}

	parent, err := GetParent(eatable)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Cannot get the parent of eatable with key " + key + " : " + err.Error()})
	}

	eatable.Parent = parent

	return context.JSON(http.StatusOK, eatable)
}

// Create creates a new eatable
// TODO we should prevent creating a new eatable with the name of an existing one
func Create(context *echo.Context) error {

	body := context.Request().Body
	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Error while reading the eatable content"})
	}

	if len(bodyBytes) == 0 {
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing eatable content"})
	}

	var eatable Eatable
	err = json.Unmarshal(bodyBytes, &eatable)
	if err != nil {
		return context.JSON(http.StatusCreated, errorMessage{"Incorrect eatable content: " + err.Error()})
	}

	// ignore the status, created_at and nutrient information that could have beeen provided
	eatable.Nutrient = nil
	eatable.Status = STATUS_NEW
	eatable.CreatedAt = date.Date{time.Now()}

	if !eatable.IsValid() {
		return context.JSON(http.StatusBadRequest, errorMessage{"Unknown type: " + eatable.Type})
	}

	updatedEatable, err := Save(&eatable)
	if err != nil {
		log.Printf("Error: cannot save eatable : %v\n", err)
		return context.JSON(http.StatusInternalServerError, errorMessage{"Eatable creation error"})
	}
	return context.JSON(http.StatusCreated, updatedEatable)
}

// Update updates an existing eatable
func Update(context *echo.Context) error {
	key := context.Param("key")
	eatable, err := FindByKey(key)
	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Cannot load eatable with key: " + key + " - " + err.Error()})
	}
	if eatable == nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Cannot load eatable with key: " + key})
	}
	body := context.Request().Body
	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Error while reading the eatable content: " + key})
	}

	//body should not be empty
	if len(bodyBytes) == 0 {
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing eatable content: " + key})
	}

	err = json.Unmarshal(bodyBytes, &eatable)
	if !eatable.IsValid() {
		return context.JSON(http.StatusBadRequest, errorMessage{"Invalid eatable type or name"})
	}

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Incorrect eatable content: " + err.Error()})
	}
	updatedEatable, err := Save(eatable)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Eatable update error"})
		log.Printf("Error: cannot update eatable: %v\n", err)
	}
	return context.JSON(http.StatusOK, updatedEatable)
}

// SetNutrient sets or modifies the nutrient information of a given eatable
func SetNutrient(context *echo.Context) error {
	key := context.Param("key")

	bodyIo := context.Request().Body
	bodyBytes, err := ioutil.ReadAll(bodyIo)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Error while reading the eatable content"})
	}

	if len(bodyBytes) == 0 {
		return context.JSON(http.StatusBadRequest, errorMessage{"Missing request body"})
	}

	eatable, errLoad := FindByKey(key)

	if errLoad != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Cannot load eatable with key: " + key + " - " + errLoad.Error()})
	}

	if eatable == nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Eatable with key " + key + " not found"})
	}

	var nutrient Nutrient
	err = json.Unmarshal(bodyBytes, &nutrient)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Incorrect nutrient content: " + err.Error()})
	}

	eatable.Nutrient = &nutrient

	eatable, err = Save(eatable)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Set Nutrient Eatable error"})
		log.Printf("Error: cannot set nutrient : %v\n", err)
	}

	return context.JSON(http.StatusOK, eatable)
}

// SetParent sets or modifies the main parent of an eatable.
func SetParent(context *echo.Context) error {
	key := context.Param("key")
	parentKey := context.Param("parentKey")

	err := SaveParent(key, parentKey)

	if err != nil {
		log.Error(context, "Impossible to set parent: "+err.Error())
		return context.JSON(http.StatusInternalServerError, errorMessage{"Impossible to set parent with key " + parentKey + " on eatble with key " + key + " - " + err.Error()})
	}

	return context.JSON(http.StatusOK, "ok")
}
