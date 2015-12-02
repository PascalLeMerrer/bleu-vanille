package eatable

// Ensures persistance of eatables in ArangoDB database

import (
	"bleuvanille/config"
	"errors"
	"fmt"

	ara "github.com/diegogub/aranGO"
)

// Save inserts an eatable into the database
func Save(eatable *Eatable) error {
	errorMap := config.Context().Save(eatable)
	if value, ok := errorMap["error"]; ok {
		return errors.New(value)
	}
	return nil
}

func SaveParent(id, parentid string) error {
	col := config.Db().Col(config.EDGENAME_EATABLE_PARENT)
	err := col.Relate("eatables/"+id, "eatables/"+parentid, map[string]interface{}{})

	if err != nil {
		return err
	}

	return nil
}

// LoadById returns the eatable object for a given id, if any
func FindById(id string) (*Eatable, error) {
	var result Eatable

	col := config.GetCollection(&result)
	err := col.Get(id, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

type Edge struct {
		Id   string `json:"_id,omitempty"  `
		From string `json:"_from"`
		To   string `json:"_to"  `	
}

type Edges struct {
	EdgesArray []Edge `json:"edges,omitempty"`
	Error bool `json:"error,omitempty"`
}

// GetParent returns the parent of a given eatable
func GetParent(child *Eatable) (*Edge, error) {
	var arrayofresult Edges
	
	col := config.Db().Col(config.EDGENAME_EATABLE_PARENT)
	err := col.Edges(child.Id, "out", &arrayofresult)

	if err != nil {
		return nil, err
	}

	if len(arrayofresult.EdgesArray) > 0 {
		return &arrayofresult.EdgesArray[0], nil
	}
	
	return nil, nil
}

// LoadByName returns the eatable object for a given name, if any
func FindByName(name string) (*Eatable, error) {
	var result Eatable

	col := config.GetCollection(&result)
	cursor, err := col.Example(map[string]interface{}{"name": name}, 0, 1)
	if err != nil {
		return nil, err
	}
	if cursor.Result != nil && len(cursor.Result) > 0 {
		cursor.FetchOne(&result)
		return &result, nil
	}
	return nil, nil
}

// LoadAll returns the list of all eatables in the database
// sort defines the sorting property name
// order must be either ASC or DESC
func FindAll(sort string, order string) ([]Eatable, error) {
	queryString := "FOR e in eatables SORT c." + sort + " " + order + " RETURN c"
	arangoQuery := ara.NewQuery(queryString)
	cursor, err := config.Db().Execute(arangoQuery)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := make([]Eatable, len(cursor.Result))
	err = cursor.FetchBatch(&result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}
