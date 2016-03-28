package statistics

import (
	"bleuvanille/config"
	"bleuvanille/log"
	"encoding/json"
	"github.com/solher/arangolite"
)

// CollectionName is the name of the collection to store metrics in the database
const CollectionName = "statistics"

// init creates the contacts collections if it do not already exist
func init() {
	config.DB().Run(&arangolite.CreateCollection{Name: CollectionName})
}

// Count retruns the value of the given counter
func Count(counterName string) (*Counter, error) {
	query := arangolite.NewQuery(` FOR counter IN %s FILTER counter.counter == @name LIMIT 1 RETURN counter `, CollectionName)
	query.Bind("name", counterName)
	rawResult, err := config.DB().Run(query)
	if err != nil {
		return nil, err
	}
	var result []Counter
	marshallErr := json.Unmarshal(rawResult, &result)
	if marshallErr != nil {
		return nil, marshallErr
	}
	if len(result) > 0 {
		return &result[0], nil
	}
	return nil, nil
}

// IncrementCounter increment the value of given counter in the database
func IncrementCounter(counterName string) {
	query := arangolite.NewQuery(` UPSERT { "counter":"%s" } INSERT { "counter":"%s", "count": 0} UPDATE { "count": OLD.count + 1 } IN %s `, counterName, counterName, CollectionName)
	_, err := config.DB().Run(query)
	if err != nil {
		log.Printf("Error: cannot increment counter %s: %v\n", counterName, err)
	}
}
