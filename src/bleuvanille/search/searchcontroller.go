package search

import (
	"errors"
	"strconv"
	
	"bleuvanille/eatablepersistance"
	"bleuvanille/log"

	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// Search searches eatable based on their name
func Search(context *echo.Context) error {
	name := context.Param("name")

	offsetParam, offsetErr := strconv.Atoi(context.Query("offset"))
	if offsetErr != nil {
		offsetParam = 0
	}
	limitParam, limitErr := strconv.Atoi(context.Query("limit"))
	if limitErr != nil {
		limitParam = 0
	}
	
	eatables,totalCount, err := SearchForEatable(name, offsetParam, limitParam)

	log.Error(context, "Error while searching for "+name)

	// Verify if the result is correctly retrieved from search
	if err != nil {
		log.Error(context, "Error while searching for "+name+" : "+err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	result := convertEatableKeyArrayInEatable(context, eatables)

	context.Response().Header().Set("X-TOTAL-COUNT", strconv.Itoa(totalCount))
	return context.JSON(http.StatusOK, result)
}

// Search searches eatable based on their name
func SearchCompletion(context *echo.Context) error {
	name := context.Param("name")

	offsetParam, offsetErr := strconv.Atoi(context.Query("offset"))
	if offsetErr != nil {
		offsetParam = 0
	}
	limitParam, limitErr := strconv.Atoi(context.Query("limit"))
	if limitErr != nil {
		limitParam = 0
	}
	
	eatables,totalCount, err := SearchPrefix(name, offsetParam, limitParam)

	log.Error(context, "Error while searching for "+name)

	// Verify if the result is correctly retrieved from search
	if err != nil {
		log.Error(context, "Error while searching for "+name+" : "+err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	result := convertEatableKeyArrayInEatable(context, eatables)

	context.Response().Header().Set("X-TOTAL-COUNT", strconv.Itoa(totalCount))
	return context.JSON(http.StatusOK, result)
}


// SearchQueryString searches eatables based on a full bleve query : used for test
func SearchQueryString(context *echo.Context) error {
	query := context.Param("query")

	offsetParam, offsetErr := strconv.Atoi(context.Query("offset"))
	if offsetErr != nil {
		offsetParam = 0
	}
	limitParam, limitErr := strconv.Atoi(context.Query("limit"))
	if limitErr != nil {
		limitParam = 0
	}

	eatables, totalCount, err := SearchFromQueryString(query, offsetParam, limitParam)

	// Verify if the result is correctly retrieved from search
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}

	result := convertEatableKeyArrayInEatable(context, eatables)

	context.Response().Header().Set("X-TOTAL-COUNT", strconv.Itoa(totalCount))
	return context.JSON(http.StatusOK, result)
}

// SearchForAllEatable returns every Eatable of the index
func SearchAllEatable(context *echo.Context) error {
	limitParam, limitErr := strconv.Atoi(context.Query("limit"))
	if limitErr != nil {
		limitParam = 0
	}

	offsetParam, offsetErr := strconv.Atoi(context.Query("offset"))
	if offsetErr != nil {
		offsetParam = 0
	}

	eatables, totalCount, err := SearchForAllEatable(offsetParam, limitParam)

	// Verify if the result is correctly retrieved from search
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}

	result := convertEatableKeyArrayInEatable(context, eatables)

	context.Response().Header().Set("X-TOTAL-COUNT", strconv.Itoa(totalCount))
	return context.JSON(http.StatusOK, result)
}

//UnIndexFromKey unindexes an eatable given its key
func UnIndexFromKey(context *echo.Context) error {
	key := context.Param("key")

	eatableVar, err := eatablepersistance.FindByKey(key)

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

//IndexFromKey indexes an eatable given its key
func IndexFromKey(context *echo.Context) error {
	key := context.Param("key")

	eatableVar, err := eatablepersistance.FindByKey(key)

	if err != nil {
		log.Error(context, "Error while indexing for "+key+" : "+err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if eatableVar == nil {
		log.Error(context, "No eatable found for key  "+key)
		return context.JSON(http.StatusNotFound, errors.New("No eatable found for key: "+key))
	}

	parent, errParent := eatablepersistance.GetParent(eatableVar)

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

//IndexAll rebuild the index from the eatable content
func IndexAll(context *echo.Context) error {
	count, err := indexAll();	

	if err != nil {
		log.Error(context, err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, count)
}

//convertEatableKeyArrayInEatable convert a list of ids to a list of real eatable struct.
func convertEatableKeyArrayInEatable(context *echo.Context, eatables []string) []eatablepersistance.Eatable {
	result := make([]eatablepersistance.Eatable, 0, len(eatables))

	for _, id := range eatables {
		//extract the key from the id
		parseId := strings.Split(id, "/")

		if len(parseId) != 2 {
			log.Error(context, "Error while retrieving the Eatable \"" + id +"\" : it has an unvalid format.")
			continue
		}

		eatableVar, err := eatablepersistance.FindByKey(parseId[1])

		if err != nil || eatableVar == nil {
			if err != nil {
				log.Error(context, "Error while retrieving the Eatable "+id+" from database : "+err.Error())
			} else {
				log.Error(context, "Error while retrieving the Eatable "+id+" from database :  eatable unknown")
			}

			continue
		}

		parent, _ := eatablepersistance.GetParent(eatableVar)
		eatableVar.Parent = parent
		result = append(result, *eatableVar)
	}

	return result
}
