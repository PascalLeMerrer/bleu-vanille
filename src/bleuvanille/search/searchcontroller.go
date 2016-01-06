package search

import (
	"errors"

	"bleuvanille/eatable"
	"bleuvanille/log"

	//	"encoding/json"
	//	"errors"
	//	"io/ioutil"
	"net/http"
	//	"time"
	"strings"

	"github.com/labstack/echo"
)

// Search searches eatable based on their name
func Search(context *echo.Context) error {
	name := context.Param("name")

	eatables, err := SearchForEatable(name)

	log.Error(context, "Error while searching for "+name)

	// Verify if the result is correctly retrieved from search
	if err != nil {
		log.Error(context, "Error while searching for "+name+" : "+err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	result := convertEatableKeyArrayInEatable(context, eatables)

	return context.JSON(http.StatusOK, result)
}

// SearchQueryString searches eatable based a full bleve query : used for test
func SearchQueryString(context *echo.Context) error {
	query := context.Param("query")

	eatables, err := SearchFromQueryString(query)

	// Verify if the result is correctly retrieved from search
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}

	result := convertEatableKeyArrayInEatable(context, eatables)

	return context.JSON(http.StatusOK, result)
}

//UnIndexFromKey unindex an eatable given its key
func UnIndexFromKey(context *echo.Context) error {
	key := context.Param("key")

	eatableVar, err := eatable.FindByKey(key)

	if err != nil {
		log.Error(context, "Error while reading eatable from  " + key + " : "+err.Error())
		log.Error(context, err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if eatableVar == nil {
		log.Error(context, "No eatable found for key  "+key)
		return context.JSON(http.StatusNotFound, errors.New("No eatable found for key: "+key))
	}

	errIndex := Delete(eatableVar)

	if errIndex != nil {
		log.Error(context, errIndex.Error())
		return context.JSON(http.StatusInternalServerError, errIndex.Error())
	}

	return context.JSON(http.StatusNoContent, nil)
}

//IndexFromKey index an eatable given its key
func IndexFromKey(context *echo.Context) error {
	key := context.Param("key")

	eatableVar, err := eatable.FindByKey(key)

	if err != nil {
		log.Error(context, "Error while indexing for "+key+" : "+err.Error())
		log.Error(context, err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if eatableVar == nil {
		log.Error(context, "No eatable found for key  "+key)
		return context.JSON(http.StatusNotFound, errors.New("No eatable found for key: "+key))
	}

	parent, errParent := eatable.GetParent(eatableVar)

	if errParent != nil {
		log.Error(context, "Error while searching for the parent of  "+key+" : "+errParent.Error())
		return context.JSON(http.StatusInternalServerError, errParent.Error())
	}

	eatableVar.Parent = parent

	errIndex := Index(eatableVar)

	if errIndex != nil {
		log.Error(context, errIndex.Error())
		return context.JSON(http.StatusInternalServerError, errIndex.Error())
	}

	return context.JSON(http.StatusOK, eatableVar)
}


//convertEatableKeyArrayInEatable convert a list of the id to a list of real eatable struct.
func convertEatableKeyArrayInEatable(context *echo.Context, eatables []string) []eatable.Eatable {
	result := make([]eatable.Eatable, len(eatables))

	var indexminus = 0

	for indexHit, id := range eatables {
		//extract the key from the id
		parseId := strings.Split(id, "/")

		if len(parseId) != 2 {
			indexminus++

			log.Error(context, "Error while retreiving the Eatable "+id+" : it has an unvalid format.")
			continue
		}

		eatableVar, err := eatable.FindByKey(parseId[1])

		if err != nil || eatableVar == nil {
			indexminus++

			if err != nil {
				log.Error(context, "Error while retreiving the Eatable "+id+" from database : "+err.Error())
			} else {
				log.Error(context, "Error while retreiving the Eatable "+id+" from database :  eatable unknown")
			}

			continue
		}

		parent, _ := eatable.GetParent(eatableVar)
		eatableVar.Parent = parent

		result[indexHit-indexminus] = *eatableVar
	}

	return result
}
