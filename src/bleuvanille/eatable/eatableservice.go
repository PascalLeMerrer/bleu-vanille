package eatable

// Ensures persistance of eatables in ArangoDB database

import (
	"bleuvanille/config"
	"encoding/json"
	"github.com/PascalLeMerrer/arangolite"
)

const CollectionName = "eatables"
const RelationshipCollectionName = "eatable_relations"

func init() {
	config.DB().Run(&arangolite.CreateCollection{Name: CollectionName})
	config.DB().Run(&arangolite.CreateCollection{Name: RelationshipCollectionName, Type: 3})
}

// Save inserts an eatable into the database
func Save(eatable *Eatable) (*Eatable, error) {
	if eatable.Key == "" {
		return create(eatable)
	} else {
		return update(eatable)
	}
}

// create inserts an eatable into the database
func create(eatable *Eatable) (*Eatable, error) {
	var query *arangolite.Query
	query = arangolite.NewQuery(` INSERT { "name": @name, "status": "%s", "type": @type, "description": @description, "createdAt": "%s" } IN %s RETURN NEW`, STATUS_NEW, eatable.CreatedAt, CollectionName)
	query.Bind("name", eatable.Name)
	query.Bind("description", eatable.Description)
	query.Bind("type", eatable.Type)
	resultByte, err := config.DB().Run(query)
	var result []Eatable
	err = json.Unmarshal(resultByte, &result)
	if err == nil {
		updatedEatable := &result[0]
		return updatedEatable, err
	}
	return nil, err
}

// update an eatable into the database
func update(eatable *Eatable) (*Eatable, error) {
	var query *arangolite.Query

	if eatable.Nutrient != nil {
		nutrient, err := json.Marshal(eatable.Nutrient)
		if err != nil {
			return nil, err
		}
		query = arangolite.NewQuery(` UPDATE "%s" WITH { "name": @name, "status": "%s", "type": @type, "description": @description, "nutrient": %s, "createdAt": "%s" } IN %s `, eatable.Key, eatable.Status, nutrient, eatable.CreatedAt.String(), CollectionName)
	} else {
		query = arangolite.NewQuery(` UPDATE "%s" WITH { "name": @name, "status": "%s", "type": @type, "description": @description, "createdAt": "%s" } IN %s `, eatable.Key, eatable.Status, eatable.CreatedAt.String(), CollectionName)
	}

	query.Bind("name", eatable.Name)
	query.Bind("description", eatable.Description)
	query.Bind("type", eatable.Type)
	_, queryErr := config.DB().Run(query)
	return eatable, queryErr
}

func SaveParent(key, parentkey string) error {
	createEdgeQuery := arangolite.NewQuery(` INSERT {'_from': @id,'_to': @parentId, 'is': 'child'} IN %s `, RelationshipCollectionName)
	createEdgeQuery.Bind("id", CollectionName+"/"+key)
	createEdgeQuery.Bind("parentId", CollectionName+"/"+parentkey)
	_, err := config.DB().Run(createEdgeQuery)
	return err
}

// FindById returns the eatable object for a given key, if any
func FindByKey(key string) (*Eatable, error) {
	return FindBy("_key", key)
}

// key GetParent returns the parent of a given eatable
func GetParent(child *Eatable) (*Eatable, error) {
	var result []Eatable

	query := arangolite.NewQuery(` FOR e IN EDGES(%s, @id, 'outbound', [ { 'is': 'child' } ]) LIMIT 1 
		For eatable in eatables FILTER eatable._id == e._to LIMIT 1 return eatable `, RelationshipCollectionName)
	query.Bind("id", child.Id)

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

// FindByName returns the eatable object for a given name, if any
func FindByName(name string) (*Eatable, error) {
	return FindBy("name", name)
}

// FindBy returns an eatable matching the given property value, or nil if not eatable matches
func FindBy(name, value string) (*Eatable, error) {
	var result []Eatable

	query := arangolite.NewQuery(` FOR e IN %s FILTER e.@name == @value LIMIT 1 RETURN e `, CollectionName)
	query.Bind("name", name)
	query.Bind("value", value)

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

// FindAll returns the list of all eatables in the database
// sort defines the sorting property name
// order must be either ASC or DESC
func FindAll(sort string, order string) ([]Eatable, error) {
	query := arangolite.NewQuery(` FOR e IN %s SORT @sort @order RETURN e `, CollectionName).Cache(true).BatchSize(500)
	query.Bind("sort", sort)
	query.Bind("order", order)
	async, asyncErr := config.DB().RunAsync(query)

	if asyncErr != nil {
		return nil, asyncErr
	}

	eatables := []Eatable{}
	decoder := json.NewDecoder(async.Buffer())

	for async.HasMore() {
		batch := []Eatable{}
		decoder.Decode(&batch)
		eatables = append(eatables, batch...)
	}

	return eatables, nil
}
