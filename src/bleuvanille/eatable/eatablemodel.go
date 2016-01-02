package eatable

import (
	"bleuvanille/date"
	"errors"
	"time"
)

const (
	TYPE_INGREDIENT_RECIPE = "ingredientrecipe"
	STATUS_NEW             = "new"
)

var validTypes = map[string]bool{
	"ingredient":        true,
	"ingredientrecette": true,
}

// Eatable is a user who registered to get information about our service
type Eatable struct {
	Id          string    `json:"_id,omitempty"`
	Key         string    `json:"_key,omitempty"`
	Parent      *Eatable  `json:"parent,omitempty"`
	Name        string    `json:"name,omitempty"`
	CreatedAt   date.Date `json:"createdAt,omitempty"`
	Status      string    `json:"status,omitempty"`
	Type        string    `json:"type,omitempty"`
	Description string    `json:"description,omitempty"`
	Nutrient    *Nutrient `json:"nutrient,omitempty"`
}

// Nutrient is the information about the composition of an eatable
type Nutrient struct {
	Carbohydrate int `json:"carbohydrate,omitempty"`
	Sugar        int `json:"sugar,omitempty"`
	Protein      int `json:"protein,omitempty"`
	Lipid        int `json:"lipid,omitempty"`
	Fiber        int `json:"fiber,omitempty"`
	Alcohol      int `json:"alcohol,omitempty"`

	// "computed" is computed from other source, "humanly-set" if a human set the value of the nutrient
	Status string `json:"status,ommitempty"`
}

// New creates an Eatable instance givent its name and description
func New(name, description string) (*Eatable, error) {

	if name == "" {
		errorMessage := "Cannot create eatable, name is missing"
		return nil, errors.New(errorMessage)
	}

	var eatable *Eatable

	eatable.Name = name
	eatable.Type = TYPE_INGREDIENT_RECIPE
	eatable.CreatedAt = date.Date{time.Now()}
	eatable.Status = STATUS_NEW
	eatable.Description = description

	return eatable, nil
}

// IsValid checks if the content of the eatable field is correct
func (eatable *Eatable) IsValid() bool {
	return eatable.Type != "" && validTypes[eatable.Type] && eatable.Name != ""
}
