package eatable

// Ensures persistance of eatables in ArangoDB database

import (
	"bleuvanille/config"
	"bleuvanille/log"
	"encoding/json"
	"github.com/PascalLeMerrer/arangolite"
)

const GraphName = "eatable_graph"
const CollectionName = "eatables"
const RelationshipCollectionName = "eatable_relations"

// init creates the eatable graph and collections if they do not already exist
func init() {
	config.DB().Run(&arangolite.CreateCollection{Name: CollectionName})
	config.DB().Run(&arangolite.CreateCollection{Name: RelationshipCollectionName, Type: 3})

	_, err := config.DB().Run(&arangolite.GetGraph{Name: GraphName})
	if err != nil {
		from := make([]string, 1)
		from[0] = CollectionName
		to := make([]string, 1)
		to[0] = CollectionName

		edgeDefinition := arangolite.EdgeDefinition{Collection: RelationshipCollectionName, From: from, To: to}
		edgeDefinitions := make([]arangolite.EdgeDefinition, 1)
		edgeDefinitions[0] = edgeDefinition
		_, err := config.DB().Run(&arangolite.CreateGraph{Name: GraphName, EdgeDefinitions: edgeDefinitions})
		if err != nil {
			log.Error(nil, "Cannot create graph with name "+GraphName+" - "+err.Error())
		}
	}
}

// Save inserts an eatable into the database
func Save(eatable *Eatable) (*Eatable, error) {
	if eatable.Key == "" {
		return insert(eatable)
	} else {
		return update(eatable)
	}
}

// insert inserts an eatable into the database
func insert(eatable *Eatable) (*Eatable, error) {

	resultByte, err := config.DB().Send("INSERT DOCUMENT", "POST", "/_api/document?collection="+CollectionName, eatable)
	err = json.Unmarshal(resultByte, eatable)
	return eatable, err
}

// update an eatable into the database
func update(eatable *Eatable) (*Eatable, error) {

	resultByte, err := config.DB().Send("UPDATE DOCUMENT", "PUT", "/_api/document/"+eatable.Id, eatable)
	err = json.Unmarshal(resultByte, eatable)
	return eatable, err
}

// SaveParent adds a relationship between two eatables
func SaveParent(key, parentkey string) error {
	createEdgeQuery := arangolite.NewQuery(` INSERT {'_from': @id,'_to': @parentId, 'is': 'child'} IN %s `, RelationshipCollectionName)
	createEdgeQuery.Bind("id", CollectionName+"/"+key)
	createEdgeQuery.Bind("parentId", CollectionName+"/"+parentkey)
	_, err := config.DB().Run(createEdgeQuery)
	return err
}

// FindByKey returns the eatable object for a given key, if any
func FindByKey(key string) (*Eatable, error) {
	return FindBy("_key", key)
}

// Remove deletes the eatable object for a given key, if any
// IT also removes relationships to or from this eatable
func Remove(key string) error {
	query := arangolite.NewQuery(`FOR e IN EDGES(%s, @id, 'any', [ { 'is': 'child' } ]) REMOVE e IN %s`, RelationshipCollectionName, RelationshipCollectionName)
	query.Bind("id", CollectionName+"/"+key)
	_, err := config.DB().Run(query)
	if err != nil {
		return err
	}
	query = arangolite.NewQuery(`REMOVE @key IN %s`, CollectionName)
	query.Bind("key", key)
	_, err = config.DB().Run(query)
	return err
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
