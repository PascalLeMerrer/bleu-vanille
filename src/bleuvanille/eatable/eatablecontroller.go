package eatable

import (
	"bleuvanille/log"

	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
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

//Return the object stored in database
func Get(context *echo.Context) error {
	id := context.Param("id")

	eatable, error := FindById(id)

	//Verify if the result is correctly retreived from database
	if error != nil {
		return context.JSON(http.StatusInternalServerError, errorMessage{"Invalid ID : " + id})
	}

	if eatable == nil {
		return context.JSON(http.StatusNotFound, errorMessage{"Invalid ID : " + id})
	}

	parentEdge, err := GetParent(eatable)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, errors.New("Impossible to get the parent of eatable "+id+" : "+err.Error()))
	}

	if parentEdge == nil {
		//The temporary struct is used to remove the fields _id, _rev and _key and add id field
		return context.JSON(http.StatusOK, struct {
			*Eatable
			Id      string `json:"id,omitempty"`
			OmitId  omit   `json:"_id,omitempty"`
			OmitRev omit   `json:"_rev,omitempty"`
			OmitKey omit   `json:"_key,omitempty"`
		}{Eatable: eatable, Id: eatable.Id})

	} else {
		result := struct {
			*Eatable
			Id       string `json:"id,omitempty"`
			ParentId string `json:"parentid"`
			OmitId   omit   `json:"_id,omitempty"`
			OmitRev  omit   `json:"_rev,omitempty"`
			OmitKey  omit   `json:"_key,omitempty"`
		}{Eatable: eatable, Id: eatable.Id, ParentId: parentEdge.To}

		return context.JSON(http.StatusOK, result)
	}
}

type omit *struct{}

// Create creates a new eatable
func Create(context *echo.Context) error {

	bodyio := context.Request().Body

	bodybytes, err := ioutil.ReadAll(bodyio)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errors.New("Error while reading the eatable content"))
	}

	if len(bodybytes) == 0 {
		return context.JSON(http.StatusBadRequest, errors.New("Missing eatable content"))
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
		return context.JSON(http.StatusCreated, errors.New("Incorrect eatable content : "+err.Error()))
	}

	//Check the type of eatable
	if eatable.Type != "" {
		existingType := validEatableType[eatable.Type]
		if existingType == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Unknown type : "+eatable.Type))
		}
	}

	//Save the new eatable object
	err = Save(eatable)

	if err != nil {
		log.Printf("Error: cannot save eatable : %v\n", err)
		return context.JSON(http.StatusInternalServerError, errors.New("Eatable creation error"))
	}

	return context.JSON(http.StatusCreated, eatable)
}

// Update updates an existing eatable
func Update(context *echo.Context) error {
	id := context.Param("id")

	//Read the body that contains the new json for our eatable
	bodyio := context.Request().Body

	bodybytes, err := ioutil.ReadAll(bodyio)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errors.New("Error while reading the eatable content : "+id))
	}

	//body should not be empty
	if len(bodybytes) == 0 {
		return context.JSON(http.StatusBadRequest, errors.New("Missing eatable content : "+id))
	}

	//Read the eatable from database
	eatable, errload := FindById(id)

	if errload != nil {
		return context.JSON(http.StatusBadRequest, errors.New("Incorrect eatable id : "+err.Error()))
	}

	//Populate the existing eatable with the new data, ignoring, status, created time and nutrient
	err = json.Unmarshal(bodybytes, &struct {
		*Eatable
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		Nutrient  *Nutrient `json:"nutrient"`
	}{Eatable: eatable})

	if err != nil {
		return context.JSON(http.StatusBadRequest, errors.New("Incorrect eatable content : "+err.Error()))
	}

	//Check the type of eatable
	if eatable.Type != "" {
		existingType := validEatableType[eatable.Type]
		if existingType == "" {
			return context.JSON(http.StatusBadRequest, errorMessage{"Unknown type : " + eatable.Type})
		}
	}

	//Save the new eatable object
	err = Save(eatable)

	if err != nil {
		log.Printf("Error: cannot update eatable : %v\n", err)
		return context.JSON(http.StatusInternalServerError, errors.New("Update eatable error"))
	}

	return context.JSON(http.StatusOK, eatable)
}

// SetNutrient sets or modifies the nutrient information of a given eatable
func SetNutrient(context *echo.Context) error {
	id := context.Param("id")

	bodyIo := context.Request().Body

	bodyBytes, err := ioutil.ReadAll(bodyIo)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errors.New("Error while reading the eatable content"))
	}

	if len(bodyBytes) == 0 {
		return context.JSON(http.StatusBadRequest, errors.New("Missing eatable content"))
	}

	eatable, errLoad := FindById(id)

	if errLoad != nil {
		return context.JSON(http.StatusBadRequest, errors.New("Incorrect eatable id : "+errLoad.Error()))
	}

	var nutrient Nutrient
	err = json.Unmarshal(bodyBytes, &nutrient)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errors.New("Incorrect nutrient content : "+err.Error()))
	}

	eatable.Nutrient = &nutrient

	err = Save(eatable)

	if err != nil {
		log.Printf("Error: cannot set nutrient : %v\n", err)
		return context.JSON(http.StatusInternalServerError, errors.New("Set Nutrient Eatable error"))
	}

	return context.JSON(http.StatusOK, eatable)
}

// SetParent sets or modifies the main parent of an eatable.
func SetParent(context *echo.Context) error {
	id := context.Param("id")
	idParent := context.Param("newParentId")

	err := SaveParent(id, idParent)

	if err != nil {
		return context.JSON(http.StatusBadRequest, errors.New("Impossible to set parent : "+err.Error()))
	}

	return context.JSON(http.StatusOK, "ok")
}
