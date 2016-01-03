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

	"github.com/labstack/echo"
)

// Search searches eatable based on their name
func Search(context *echo.Context) error {
	name := context.Param("name")

	eatables, err := SearchForEatable(name)

	// Verify if the result is correctly retrieved from search
	if err != nil {
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

func IndexFromKey(context *echo.Context) error {
	key := context.Param("key")

	eatableVar, err := eatable.FindByKey(key)

	if err != nil {
		log.Error(context, err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if eatableVar == nil {
		log.Error(context, "No eatable found for key  " + key)
		return context.JSON(http.StatusNotFound, errors.New("No eatable found for key: " + key))
	}
	
	parent, errParent := eatable.GetParent(eatableVar)

	if errParent != nil {
		log.Error(context, "Error while searching for the parent of  " + key + " : " +  errParent.Error())
		return context.JSON(http.StatusInternalServerError, errParent.Error())
	}

	eatableVar.Parent = parent

	errIndex := Index(eatableVar)

	if errIndex != nil {
		log.Error(context, errIndex.Error())
		return context.JSON(http.StatusInternalServerError, errIndex.Error())
	}

	return context.JSON(http.StatusInternalServerError, eatableVar)
}

func convertEatableKeyArrayInEatable(context *echo.Context, eatables []string) []eatable.Eatable {
	result := make([]eatable.Eatable, len(eatables))

	var indexminus = 0

	for indexHit, key := range eatables {
		eatableVar, err := eatable.FindByKey(key)

		if err != nil {
			indexminus++
			log.Error(context, "Error while retreiving the Eatable " + key + " from database : "+err.Error())
		} else {
			parent, _ := eatable.GetParent(eatableVar)
			eatableVar.Parent = parent
			
			result[indexHit-indexminus] = *eatableVar
		}
	}

	return result
}
