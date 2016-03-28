package statistics

import (
	"github.com/labstack/echo"
	"net/http"
)

// Get returns the entry for a given email if any
func Get() echo.HandlerFunc {
	return func(context echo.Context) error {
		counterName := context.Param("counter")
		counter, err := Count(counterName)
		if err == nil {
			return context.JSON(http.StatusOK, counter)
		}
		return context.String(http.StatusInternalServerError, err.Error())
	}
}
