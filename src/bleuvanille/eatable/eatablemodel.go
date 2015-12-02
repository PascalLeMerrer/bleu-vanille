package eatable

import (
	"bleuvanille/config"
	"errors"
	"time"

	ara "github.com/diegogub/aranGO"
)

const (
	TYPE_INGREDIENT_RECIPE = "ingredientrecipe"
	STATUS_NEW            = "new"
)

// Eatable is a user who registered to get information about our service
type Eatable struct {
	ara.Document
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Nutrient    *Nutrient `json:"nutrient,omitempty"`
}

// Nutrient is the information about the composition of an eatable
type Nutrient struct {
	Carbohydrate int `json:"carbohydrate"`
	Sugar        int `json:"sugar"`
	Protein      int `json:"protein"`
	Lipid        int `json:"lipid"`
	Fiber        int `json:"fiber"`
	Alcohol      int `json:"alcohol"`
	//       int    `json:""`

	// "computed" is computed from other source, "humanly-set" if a human set the value of the nutrient
	Status string `json:"status,ommitempty"`
}

// Eatables is a list of Eatable
type Eatables []Eatable

// NewEmpty creates an Eatable instance with the default value
func NewEmpty() *Eatable {
	var eatable Eatable

	eatable.Name = ""
	eatable.CreatedAt = time.Now()
	eatable.Description = ""
	eatable.Type = TYPE_INGREDIENT_RECIPE
	eatable.Status = STATUS_NEW

	return &eatable
}

// New creates an Eatable instance givent its name and description
func New(name, description string) (*Eatable, error) {

	if name == "" {
		errorMessage := "Cannot create eatable, name is missing"
		return nil, errors.New(errorMessage)
	}

	var eatable *Eatable

	eatable = NewEmpty()
	eatable.Name = name
	eatable.Description = description

	return eatable, nil
}

// GetKey returns the key in ArangoDB for the eatable
func (eatable *Eatable) GetKey() string {
	return eatable.Key
}

// GetCollection returns the collection name in ArangoDB for contacts
func (eatable *Eatable) GetCollection() string {
	return config.COLNAME_ETABLES
}

// GetError returns true if there is an error and gives the last error message
func (eatable *Eatable) GetError() (string, bool) {
	return eatable.Message, eatable.Error
}
