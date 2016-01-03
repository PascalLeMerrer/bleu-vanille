package search

import (
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

	result := convertEatableInExportable(context, eatables)

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

	result := convertEatableInExportable(context, eatables)

	return context.JSON(http.StatusOK, result)
}

func IndexFromId(context *echo.Context) error {
	id := context.Param("id")

	exportableEatable, err := eatable.GetExportableEatable(id)

	if exportableEatable.ParentId != "" {
		parent, errParent := eatable.FindById(exportableEatable.ParentId)

		if errParent != nil {
			log.Error(context, "Impossible to read the parent of eatable  "+exportableEatable.Id+" : "+errParent.Error())
		} else {
			if parent != nil {
				exportableEatable.ParentName = parent.Name
			} else {
				log.Error(context, "Impossible to read the parent of eatable  "+exportableEatable.Id)
			}
		}
	}

	if err != nil {
		log.Error(context, err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	errIndex := Index(exportableEatable)

	if errIndex != nil {
		log.Error(context, errIndex.Error())
		return context.JSON(http.StatusInternalServerError, errIndex.Error())
	}

	return context.JSON(http.StatusInternalServerError, exportableEatable)
}

//	echoServer.Get("/search/query/:query", search.SearchQueryString)
//	echoServer.Get("/search/index/:id", search.IndexFromId)

func convertEatableInExportable(context *echo.Context, eatables []string) []eatable.ExportableEatable {
	result := make([]eatable.ExportableEatable, len(eatables))

	var indexminus = 0

	for indexHit, id := range eatables {
		exportableEatable, err := eatable.GetExportableEatable(id)

		if err != nil {
			indexminus++
			log.Error(context, "Error while retreiving the Eatable "+id+" from database : "+err.Error())
		} else {
			result[indexHit-indexminus] = *exportableEatable
		}
	}

	return result
}
