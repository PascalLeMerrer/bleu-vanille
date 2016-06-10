package ingredient

import (
	"bleuvanille/session"
	"errors"
	"fmt"
	"github.com/goodsign/monday"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type formattedIngredient struct {
	ID          string `json:"_id"`
	Key         string `json:"_key,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt"`
	Creator     string `json:"creator"`
	// Energy is defined in KJ/100g
	Energy   int    `json:"energy,omitempty"`
	Category string `json:"category,omitempty"`
	// False when the ingredient was created by an user, and needs to be verified by an admin
	Approved bool `json:"approved"`
	// the number of the months when this ingredient is naturally available
	Months []int `json:"months,omitempty"`
}

// Get returns an ingredient description
func Get() echo.HandlerFunc {
	return func(context echo.Context) error {
		ingredient, err := FindByKey(context.Param("key"))
		if err != nil || ingredient == nil {
			log.Println(err)
			return context.JSON(http.StatusNotFound, nil)
		}
		return context.JSON(http.StatusOK, formatIngredient(ingredient))
	}
}

// GetAll writes the list of all ingredients
func GetAll() echo.HandlerFunc {
	return func(context echo.Context) error {
		offsetParam, offsetErr := strconv.Atoi(context.QueryParam("offset"))
		if offsetErr != nil {
			offsetParam = 0
		}
		limitParam, limitErr := strconv.Atoi(context.QueryParam("limit"))
		if limitErr != nil {
			limitParam = 0
		}

		criteria, order := getSortingParams(context)
		ingredients, totalCount, err := FindAll(criteria, order, offsetParam, limitParam)

		if err != nil {
			return context.JSON(http.StatusInternalServerError, errors.New("Ingredient list retrieval error"))
		}
		formattedIngredients := make([]formattedIngredient, len(ingredients))
		for i := range ingredients {
			formattedIngredients[i] = formatIngredient(&ingredients[i])
			i++
		}
		context.Response().Header().Set("X-TOTAL-COUNT", strconv.Itoa(totalCount))
		return context.JSON(http.StatusOK, formattedIngredients)
	}
}

// converts an Ingredient object to a FormattedIngredient object
// which is its public counterpart
func formatIngredient(ingredient *Ingredient) formattedIngredient {
	formattedDate := monday.Format(ingredient.CreatedAt, "Mon _2 Jan 2006 15:04", monday.LocaleFrFR)
	return formattedIngredient{ingredient.ID, ingredient.Key, ingredient.Name, ingredient.Description, formattedDate, ingredient.Creator, ingredient.Energy, ingredient.Category, ingredient.Approved, ingredient.Months}
}

// extracts from the sortParam parameter a criteria value and a sorting order
// @returns the nma of the property on which user list must be sorted, and the sort order (ASC or DESC)
func getSortingParams(context echo.Context) (string, string) {

	sortParam := context.QueryParam("sort")

	var criteria string
	var order string

	switch sortParam {
	case "newer":
		criteria = "createdAt"
		order = "DESC"
	case "older":
		criteria = "createdAt"
		order = "ASC"
	case "nameAsc":
		criteria = "lastname"
		order = "ASC"
	case "nameDesc":
		criteria = "lastname"
		order = "DESC"
	default:
		criteria = "createdAt"
		order = "DESC"
	}
	return criteria, order
}

// Create creates a new ingredient
func Create() echo.HandlerFunc {
	return func(context echo.Context) error {
		ingredient := &Ingredient{}
		session := context.Get("session").(*session.Session)
		ingredient.Creator = session.UserID
		if err := context.Bind(ingredient); err != nil {
			return context.JSON(http.StatusBadRequest, errors.New("Invalid JSON parameter to create an ingredient"))
		}

		result, err := Save(ingredient)
		if err != nil {
			if err.Error() == "cannot create document, unique constraint violated" {
				return context.JSON(http.StatusConflict, errors.New("Ingredient is already registered"))
			}
			log.Printf("Error: cannot save ingredient: %+v: %v\n", ingredient, err)
			return context.JSON(http.StatusInternalServerError, errors.New("Ingredient creation error"))
		}

		return context.JSON(http.StatusCreated, result)
	}
}

// Patch modifies a  given ingredient
func Patch() echo.HandlerFunc {
	return func(context echo.Context) error {
		ingredient, err := FindByKey(context.Param("key"))
		if err != nil || ingredient == nil {
			log.Println(err)
			return context.JSON(http.StatusNotFound, nil)
		}

		status := ingredient.Approved
		if err := context.Bind(ingredient); err != nil {
			log.Printf("Cannot bind JSON %+v for ingredient update: %s \n", context.Request(), err)
			return context.JSON(http.StatusBadRequest, errors.New("Invalid JSON parameter for ingredient update"))
		}
		session := context.Get("session").(*session.Session)

		if (status != ingredient.Approved) && !session.IsAdmin {
			return context.JSON(http.StatusForbidden, errors.New("User not allowed to modify ingredient status"))
		}
		result, err := Save(ingredient)
		if err != nil {
			log.Printf("Error: cannot save ingredient: %+v: %v\n", ingredient, err)
			return context.JSON(http.StatusInternalServerError, errors.New("Ingredient update error"))
		}
		return context.JSON(http.StatusOK, formatIngredient(result))
	}
}

// Delete removes from the database the ingredient for a given key
// the ingredient musn't be approved
// TODO if the user is not an admin, it can delete ingredients he created
func Delete() echo.HandlerFunc {
	return func(context echo.Context) error {
		key := context.Param("key")

		fmt.Printf("DELETE %s \n", key)
		if key == "" {
			return context.JSON(http.StatusBadRequest, errors.New("Missing key parameter in DELETE request"))
		}

		ingredient, err := FindByKey(context.Param("key"))
		if err != nil || ingredient == nil {
			log.Println(err)
			return context.JSON(http.StatusNotFound, nil)
		}

		if ingredient.Approved {
			return context.JSON(http.StatusForbidden, "An approved ingredient cannot be deleted")
		}

		session := context.Get("session").(*session.Session)
		if !session.IsAdmin && ingredient.Creator != session.UserID {
			return context.JSON(http.StatusForbidden, "Insufficient rights to delete this ingredient")
		}

		err = Remove(key)

		if err != nil {
			log.Printf("Cannot delete ingredient with key %s, error: %s", key, err)
			return context.JSON(http.StatusInternalServerError, fmt.Errorf("Cannot delete ingredient with key %s", key))
		}
		return context.NoContent(http.StatusNoContent)
	}
}
