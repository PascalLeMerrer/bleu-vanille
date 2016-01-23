package eatable

import (
	"bleuvanille/date"
	"bleuvanille/log"
	"bleuvanille/eatablepersistance"
	"bleuvanille/search"
		
	"encoding/json"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"time"
	"errors"
)

type errorMessage struct {
	Message string `json:"error"`
}

// Get returns the eatable object stored in database
func Get(context *echo.Context) error {
	key := context.Param("key")

	eatable, error := eatablepersistance.FindByKey(key)

	if error != nil {
		log.Error(context, error.Error())
		return context.JSON(http.StatusInternalServerError, errorMessage{"Invalid key: " + key})
	}

	if eatable == nil {
		return context.JSON(http.StatusNotFound, errorMessage{"No eatable found for key: " + key})
	}

	parent, err := eatablepersistance.GetParent(eatable)

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

	var eatable eatablepersistance.Eatable
	err = json.Unmarshal(bodyBytes, &eatable)
	if err != nil {
		return context.JSON(http.StatusCreated, errorMessage{"Incorrect eatable content: " + err.Error()})
	}

	// ignore the status, created_at and nutrient information that could have beeen provided
	eatable.Nutrient = nil
	eatable.Status = eatablepersistance.STATUS_NEW
	eatable.CreatedAt = date.Date{time.Now()}

	if !eatable.IsValid() {
		return context.JSON(http.StatusBadRequest, errorMessage{"Unknown type: " + eatable.Type})
	}

	existingEatable, _ := eatablepersistance.FindByName(eatable.Name)
	if existingEatable != nil {
		return context.JSON(http.StatusConflict, errorMessage{"Eatable with name " + eatable.Name + " already exists."})
	}

	updatedEatable, err := eatablepersistance.Save(&eatable)
	if err != nil {
		log.Printf("Error: cannot save eatable : %v\n", err)
		return context.JSON(http.StatusInternalServerError, errorMessage{"Eatable creation error"})
	}

	//Search Indexation
	err = search.Index(&eatable)

	if err != nil {
		log.Error(nil, "Error while updating the index for the eatable "+ eatable.Id +" : "+err.Error())
	}
	
	return context.JSON(http.StatusCreated, updatedEatable)
}

// Delete removes an existing eatable from the database
// Mainly intended for removing test data
// For real eatables you should prefer turning their status to disabled
func Delete(context *echo.Context) error {
	key := context.Param("key")
		
	err := eatablepersistance.Remove(key)
	if err == nil {
		//Search Delete
		err = search.DeleteFromId("eatables/" + key)
		
		if err != nil {
			log.Error(nil, "Error while desindexing the eatable "+ key +" : "+err.Error())
		}
			
		return context.String(http.StatusNoContent, "")
	} else {
		return context.JSON(http.StatusForbidden, errorMessage{"Cannot remove eatable with key " + key})
	}
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
	updatedEatable, err := eatablepersistance.Save(eatable)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Eatable update error"})
		log.Printf("Error: cannot update eatable: %v\n", err)
	}
	
	//Search Indexation
	err = search.Index(eatable)

	if err != nil {
		log.Error(nil, "Error while updating the index for the eatable "+ eatable.Id +" : "+err.Error())
	}

	return context.JSON(http.StatusOK, updatedEatable)
}

// Patch modifies the user account for a given ID
// This is an admin feature, not supposed to be used by normal users
func Patch(context *echo.Context) error {
	key := context.Param("key")

	eatable, err := eatablepersistance.FindByKey(key)

	if err != nil {
		log.Error(context, err.Error())
		return context.JSON(http.StatusInternalServerError, errorMessage{"Invalid key: " + key})
	}

	if eatable == nil {
		return context.JSON(http.StatusNotFound, errorMessage{"No eatable found for key: " + key})
	}

	err = context.Bind(eatable)
	if err != nil {
		log.Printf("Cannot bind eatable %v", err)
		return context.JSON(http.StatusBadRequest, errors.New("Cannot decode request body"))
	}

	updatedEatable, err := eatablepersistance.Save(eatable)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Eatable update error"})
		log.Printf("Error: cannot update eatable: %v\n", err)
	}
	
	//Search Indexation
	err = search.Index(eatable)

	if err != nil {
		log.Error(nil, "Error while updating the index for the eatable "+ eatable.Id +" : "+err.Error())
	}

	return context.JSON(http.StatusOK, updatedEatable)	
}

// SetStatus sets or modifies the nutrient information of a given eatable
func SetStatus(context *echo.Context) error {
	eatable, bodyBytes, errMessage := prepareUpdate(context)
	if errMessage != "" {
		return context.JSON(http.StatusBadRequest, errorMessage{errMessage})
	}

	var tempEatable eatablepersistance.Eatable
	err := json.Unmarshal(bodyBytes, &tempEatable)
	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Incorrect nutrient content: " + err.Error()})
	}

	if !eatable.IsValid() {
		return context.JSON(http.StatusBadRequest, errorMessage{"Invalid eatable type"})
	}

	eatable.Status = tempEatable.Status

	eatable, err = eatablepersistance.Save(eatable)

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

	var nutrient eatablepersistance.Nutrient
	err := json.Unmarshal(bodyBytes, &nutrient)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errorMessage{"Incorrect nutrient content: " + err.Error()})
	}

	eatable.Nutrient = &nutrient

	eatable, err = eatablepersistance.Save(eatable)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Set Eatable Nutrient error"})
		log.Printf("Error: cannot set nutrient : %v\n", err)
	}

	return context.JSON(http.StatusOK, eatable)
}

// loads the eatable matching the key param of the request
// and decode the request body content to be applied to this eatable
// returns the eatable, the property to be modified and an error message or an empty string
func prepareUpdate(context *echo.Context) (*eatablepersistance.Eatable, []byte, string) {

	bodyIo := context.Request().Body
	bodyBytes, err := ioutil.ReadAll(bodyIo)

	if err != nil {
		return nil, nil, "Error while reading the body content"
	}

	if len(bodyBytes) == 0 {
		return nil, nil, "Missing request body"
	}

	key := context.Param("key")

	eatable, errLoad := eatablepersistance.FindByKey(key)

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

	err := eatablepersistance.SaveParent(key, parentKey)

	if err != nil {
		log.Error(context, "Impossible to set parent: "+err.Error())
		return context.JSON(http.StatusInternalServerError, errorMessage{"Impossible to set parent with key " + parentKey + " on eatble with key " + key + " - " + err.Error()})
	}

	//Search Indexation
	eatable, err := eatablepersistance.FindByKey(key)

	if err != nil {
		log.Error(context, err.Error())
		return context.JSON(http.StatusInternalServerError, errorMessage{"Invalid key: " + key})
	}

	if eatable == nil {
		return context.JSON(http.StatusNotFound, errorMessage{"No eatable found for key: " + key})
	}

	err = search.Index(eatable)

	return context.JSON(http.StatusOK, "ok")
}
