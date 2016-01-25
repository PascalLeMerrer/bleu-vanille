package search

import (
	"fmt"
	"testing"

	"bleuvanille/eatable"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndex(t *testing.T) {
	index, err := getIndex()

	assert.NoError(t, err, "Error when creating the index")
	assert.NotNil(t, index, "Error when creating the index : nil index")
}

func TestBasicSearch(t *testing.T) {

	eatable1 := eatable.Eatable{
		Name:        "Légume",
		Description: "titi",
		Type:        "Ingredient",
	}

	eatable2 := eatable.Eatable{
		Name:        "Carotte",
		Description: "Description2",
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
	eatable.SaveParent(eatable2.Key, eatable1.Key)

	index, errCreationIndex := getIndex()

	assert.NoError(t, errCreationIndex, "Error when creating the index")
	assert.NotNil(t, index, "Error when creating the index")

	parent1, errParent1 := eatable.GetParent(&eatable1)
	eatable1.Parent = parent1
	assert.NoError(t, errParent1, "Error when fetching eatable1")
	errIndex1 := Index(&eatable1)
	assert.NoError(t, errIndex1, "Error when indexing eatable1")
	defer Delete(&eatable1)

	parent2, errParent2 := eatable.GetParent(&eatable2)
	eatable2.Parent = parent2

	assert.NoError(t, errParent2, "Error when fetching eatable2")
	errIndex2 := Index(&eatable2)
	assert.NoError(t, errIndex2, "Error when indexing eatable2")
	defer Delete(&eatable2)

	parent3, errParent3 := eatable.GetParent(&eatable3)
	eatable3.Parent = parent3
	assert.NoError(t, errParent3, "Error when fetching eatable3")
	errIndex3 := Index(&eatable3)
	assert.NoError(t, errIndex3, "Error when indexing eatable3")
	defer Delete(&eatable3)

	//Search for the main ingredient : find one
	{
		results, errSearch := SearchForEatable("Carotte")
		assert.NoError(t, errSearch, "Error when searching for Carotte")
		assert.True(t, len(results) == 1, "Error when searching for Carotte : %d results returned instead of 1", len(results))

		if len(results) == 1 {
			assert.True(t, results[0] == eatable2.Id, "Error when searching for Carotte : Id is not correct. Should be %s but is %s", eatable2.Id, results[0])
		}
	}

	//Search for the parent ingredient : find two
	{
		results, errSearch := SearchForEatable("Légume")
		assert.NoError(t, errSearch, "Error when searching for Légume")
		assert.True(t, len(results) == 2, "Error when searching for Légume : %d results returned instead of 2", len(results))
	}

	//Search for the parent ingredient without accent : find two
	{
		results, errSearch := SearchForEatable("Legume")
		assert.NoError(t, errSearch, "Error when searching for Legume")
		assert.True(t, len(results) == 2, "Error when searching for Legume : %d results returned instead of 2", len(results))
	}

	//Search for carote with a spelling mistake : find two
	{
		results, errSearch := SearchForEatable("carote")
		assert.NoError(t, errSearch, "Error when searching for carote")
		assert.True(t, len(results) == 1, "Error when searching for carote : %d results returned instead of 1", len(results))
	}
}

