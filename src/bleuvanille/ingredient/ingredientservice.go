package ingredient

// Ensures persistance of ingredients in Postgresql database

import (
	"bleuvanille/config"
	"encoding/json"
	"fmt"
	"github.com/solher/arangolite"
)

// CollectionName is the name of the ingredient collection in the database
const CollectionName = "ingredients"

// init creates the ingredients collections if it do not already exist
func init() {
	indexFields := []string{"name"}
	config.CreateHashIndexedCollection(CollectionName, indexFields)
}

// Save inserts an ingredient into the database
func Save(ingredient *Ingredient) (*Ingredient, error) {

	// ingredientJSON, err := json.Marshal(ingredient)
	// if err != nil {
	// 	return nil, err
	// }

	//query := arangolite.NewQuery(` UPSERT { "name":"%s" } INSERT %s REPLACE %s IN %s RETURN NEW`, ingredient.Name, ingredientJSON, ingredientJSON, CollectionName)
	var resultByte []byte
	var err error

	if ingredient.Key == "" || ingredient.ID == "" {
		resultByte, err = config.DB().Send("INSERT DOCUMENT", "POST", "/_api/document?collection="+CollectionName, ingredient)
	} else {
		query := fmt.Sprintf("/_api/document/%s/%s", CollectionName, ingredient.Key)
		resultByte, err = config.DB().Send("UPDATE DOCUMENT", "PUT", query, ingredient)
	}

	// resultByte, err := config.DB().Run(query)
	// var result []Ingredient
	if err == nil {
		// err = json.Unmarshal(resultByte, &result)
		err = json.Unmarshal(resultByte, ingredient)
	}
	// if err == nil {
	// 	ingredient = &result[0]
	// }

	return ingredient, err
}

// FindByKey returns the ingredient object for a given ingredient key, if any
func FindByKey(key string) (*Ingredient, error) {
	return FindBy("_key", key)
}

// FindByName returns the ingredient object for a given name, if any
func FindByName(name string) (*Ingredient, error) {
	return FindBy("name", name)
}

// FindBy returns an ingredient matching the given property value, or nil if not ingredient matches
func FindBy(name, value string) (*Ingredient, error) {

	query := arangolite.NewQuery(` FOR u IN %s FILTER u.@name == @value LIMIT 1 RETURN u `, CollectionName)
	query.Bind("name", name)
	query.Bind("value", value)

	return executeReadingQuery(query)
}

// Executes a given query that is expected to return a single ingredient
func executeReadingQuery(query *arangolite.Query) (*Ingredient, error) {
	var result []Ingredient

	rawResult, err := config.DB().Run(query)
	if err != nil {
		return nil, err
	}

	marshallErr := json.Unmarshal(rawResult, &result)
	if marshallErr != nil {
		return nil, marshallErr
	}
	if len(result) > 0 {
		return &result[0], nil
	}
	return nil, nil
}

// FindAll returns the list of all ingredients in the database
// sort defines the sorting property name
// order must be either ASC or DESC
// offset is the start index
// limit défines the max number of results to be returned
// returns an array of ingredients, the total number of ingredients, an error
func FindAll(sort string, order string, offset int, limit int) ([]Ingredient, int, error) {
	fmt.Printf("\n\n<<<<< FindAll >>>>>> sort %s, order %s, offset %v, limit %v", sort, order, offset, limit)

	limitString := ""
	if limit > 0 {
		limitString = fmt.Sprintf("LIMIT %d, %d", offset, limit)
	}
	queryString := "FOR ingredient IN %s SORT ingredient.%s %s %s RETURN ingredient"
	query := arangolite.NewQuery(queryString, CollectionName, sort, order, limitString).BatchSize(100)

	async, asyncErr := config.DB().RunAsync(query)
	if asyncErr != nil {
		return nil, 0, asyncErr
	}

	ingredients := []Ingredient{}
	decoder := json.NewDecoder(async.Buffer())

	for async.HasMore() {
		batch := []Ingredient{}
		decoder.Decode(&batch)
		ingredients = append(ingredients, batch...)
	}

	return ingredients, len(ingredients), nil
}

// Remove removes the entry for a given key
func Remove(key string) error {
	query := arangolite.NewQuery(`FOR ingredient IN %s FILTER ingredient._key==@key REMOVE ingredient IN %s`, CollectionName, CollectionName)
	query.Bind("key", key)
	_, err := config.DB().Run(query)
	return err
}
