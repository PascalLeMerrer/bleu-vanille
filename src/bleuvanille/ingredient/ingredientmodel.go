package ingredient

import (
	"errors"
	"log"
	"time"
)

// Ingredient is basic component for a recipe
type Ingredient struct {
	ID          string    `json:"_id,omitempty"`
	Key         string    `json:"_key,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	// Energy is defined in KJ/100g
	Energy   int    `json:"energy,omitempty"`
	Category string `json:"category,omitempty"`
	// False when the ingredient was created by an user, and needs to be verified by an admin
	Approved bool `json:"approved"`
	// the ID of the user that created this ingredient
	Creator string `json:"creator"`
	// the number of the months when this ingredient is naturally available
	Months []int `json:"months,omitempty"`
}

// New creates a Ingredient instance
// months are the number of the months when this ingredient is naturally available
// 1 = January
func New(ID string, name string, description string, energy int, category string, months []int, creator string) (Ingredient, error) {
	var ingredient Ingredient
	ingredient.ID = ID
	if name == "" {
		errorMessage := "Cannot create ingredient, name is missing"
		log.Println(errorMessage)
		return ingredient, errors.New(errorMessage)
	}
	ingredient.Name = name
	ingredient.Description = description
	ingredient.CreatedAt = time.Now()
	ingredient.Energy = energy
	ingredient.Category = category
	ingredient.Months = months
	ingredient.Creator = creator
	return ingredient, nil
}
