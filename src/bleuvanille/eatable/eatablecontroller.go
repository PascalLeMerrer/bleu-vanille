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

	existingEatable, _ := FindByName(eatable.Name)
	if existingEatable != nil {
		return context.JSON(http.StatusConflict, errorMessage{"Eatable with name " + eatable.Name + " already exists."})
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
	eatable, bodyBytes, errMessage := prepareUpdate(context)
	if errMessage != "" {
		return context.JSON(http.StatusBadRequest, errorMessage{errMessage})
	}

	err := json.Unmarshal(bodyBytes, &eatable)
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

// SetStatus sets or modifies the nutrient information of a given eatable
func SetStatus(context *echo.Context) error {
	eatable, bodyBytes, errMessage := prepareUpdate(context)
	if errMessage != "" {
		return context.JSON(http.StatusBadRequest, errorMessage{errMessage})
	}

	var tempEatable Eatable
	err := json.Unmarshal(bodyBytes, &tempEatable)
	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Incorrect nutrient content: " + err.Error()})
	}

	if !eatable.IsValid() {
		return context.JSON(http.StatusBadRequest, errorMessage{"Invalid eatable type"})
	}

	eatable.Status = tempEatable.Status

	eatable, err = Save(eatable)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Set Eatable Status error"})
		log.Printf("Error: cannot set status : %v\n", err)
	}

	return context.JSON(http.StatusOK, eatable)
}

// SetNutrient sets or modifies the nutrient information of a given eatable
func SetNutrient(context *echo.Context) error {
	eatable, bodyBytes, errMessage := prepareUpdate(context)

	if errMessage != "" {
		return context.JSON(http.StatusBadRequest, errorMessage{errMessage})
	}

	var nutrient Nutrient
	err := json.Unmarshal(bodyBytes, &nutrient)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Incorrect nutrient content: " + err.Error()})
	}

	eatable.Nutrient = &nutrient

	eatable, err = Save(eatable)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Set Eatable Nutrient error"})
		log.Printf("Error: cannot set nutrient : %v\n", err)
	}

	return context.JSON(http.StatusOK, eatable)
}

// loads the eatable matching the key param of the request
// and decode the request body content to be applied to this eatable
// returns the eatable, the property to be modified and an error message or an empty string
func prepareUpdate(context *echo.Context) (*Eatable, []byte, string) {

	bodyIo := context.Request().Body
	bodyBytes, err := ioutil.ReadAll(bodyIo)

	if err != nil {
		return nil, nil, "Error while reading the body content"
	}

	if len(bodyBytes) == 0 {
		return nil, nil, "Missing request body"
	}

	key := context.Param("key")

	eatable, errLoad := FindByKey(key)

	if errLoad != nil {
		return nil, nil, "Cannot load eatable with key: " + key + " - " + errLoad.Error()
	}

	if eatable == nil {
		return nil, nil, "Eatable with key " + key + " not found"
	}

	return eatable, bodyBytes, ""
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
