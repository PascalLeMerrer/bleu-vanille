package search

import (
	"testing"
	
	"bleuvanille/eatable"
	"bleuvanille/log"

	"github.com/stretchr/testify/assert"
	//	 "github.com/blevesearch/bleve"
)

//func TestCreateFakeIndex(t *testing.T) {
//	mapping := bleve.NewIndexMapping()
//    index, err := bleve.New("example.bleve", mapping)
//    if err != nil {
//        return
//    }
//
//    if index != nil {
//    	return
//    }
//}

func TestCreateIndex(t *testing.T) {
	index, err := getIndex()

	assert.NoError(t, err, "Error when creating the index")
	assert.NotNil(t, index, "Error when creating the index : nil index")
}

func TestBasicSearch(t *testing.T) {

	eatable1 := eatable.Eatable{
		Name:        "Légutestcsame",
		Description: "titi",
		Type:        "Ingredient",
	}

	eatable2 := eatable.Eatable{
		Name:        "Catestcsarotte",
		Description: "testcsa2",
		Type:        "Ingredient",
	}

	eatable3 := eatable.Eatable{
		Name:        "ingredientfake",
		Description: "",
		Type:        "Ingredient",
	}

	eatable.Save(&eatable1)
	eatable.Save(&eatable2)
	eatable.Save(&eatable3)
	eatable.SaveParent(eatable2.GetKey(), eatable1.GetKey())

	index, errCreationIndex := getIndex()

	assert.NoError(t, errCreationIndex, "Error when creating the index")
	assert.NotNil(t, index, "Error when creating the index")

	ee1, errFetch1 := eatable.GetExportableEatable(eatable1.GetKey())
	assert.NoError(t, errFetch1, "Error when fetching eatable1")
	errIndex1 := Index(ee1)
	assert.NoError(t, errIndex1, "Error when indexing eatable1")
	defer Delete(ee1)

	ee2, errFetch2 := eatable.GetExportableEatable(eatable2.GetKey())
	log.Error(nil, ee2.ParentName)
		

	assert.NoError(t, errFetch2, "Error when fetching eatable2")
	errIndex2 := Index(ee2)
	assert.NoError(t, errIndex2, "Error when indexing eatable2")
//	defer Delete(ee2)

	ee3, errFetch3 := eatable.GetExportableEatable(eatable3.GetKey())
	assert.NoError(t, errFetch3, "Error when fetching eatable3")
	errIndex3 := Index(ee3)
	assert.NoError(t, errIndex3, "Error when indexing eatable3")
	defer Delete(ee3)

	results, errSearch := SearchForEatable("Catestcsarotte")
	assert.NoError(t, errSearch, "Error when searching for Catestcsarotte")
	assert.True(t, len(results) == 1, "Error when searching for Catestcsarotte : %d results returned", len(results))
	assert.True(t, results[0] == eatable2.Id, "Error when searching for Catestcsarotte : Id is not correct. Should be %s", eatable2.Id)

}
