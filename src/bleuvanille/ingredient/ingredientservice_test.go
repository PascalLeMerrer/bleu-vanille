package ingredient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	name := "ingredient_TestCreate"
	key := "90101"
	err := Remove(key)
	assert.Nil(t, err, "Deleting an ingredient should not return an error")
	months := []int{4, 5, 6}

	ingredient := Ingredient{Key: key, Name: name, Months: months, Approved: true, Energy: 12, Category: "anything", Description: "un texte"}
	result, err := Save(&ingredient)
	assert.NoError(t, err, "Saving a new ingredient should not fail")
	assert.NotNil(t, result, "Saving a new ingredient should return it")
	assert.Equal(t, result.Key, ingredient.Key)
	assert.Equal(t, result.Name, ingredient.Name)
	assert.Equal(t, result.Approved, ingredient.Approved)
	assert.Equal(t, result.Energy, ingredient.Energy)
	assert.Equal(t, result.Months, ingredient.Months)
	err = Remove(key)
	assert.Nil(t, err, "Deleting an ingredient should not return an error")
}

func TestDuplicateName(t *testing.T) {
	var name = "ingredient TestDuplicateName"
	key1 := "90102"
	key2 := "90103"
	ingredient := Ingredient{Key: key1, Name: name}
	result, err := Save(&ingredient)
	assert.NoError(t, err, "Saving ingredient")
	assert.Equal(t, result.Key, ingredient.Key)
	assert.Equal(t, result.Name, ingredient.Name)

	ingredient2 := Ingredient{Key: key2, Name: name}
	_, err = Save(&ingredient2)
	assert.Error(t, err, "Saving ingredient with duplicate name should fail")

	err = Remove(key1)
	assert.Nil(t, err, "Deleting ingredient should not return an error")
	err = Remove(key2)
	assert.Nil(t, err, "Deleting ingredient should not return an error")
}

func TestDeleteUnknownName(t *testing.T) {
	err := Remove("99999")
	assert.Nil(t, err, "Deleting non existant name should not return an error")
}

func TestSearchByName(t *testing.T) {
	name := "TestSearchByName"
	key := "90104"
	//Ensure dat of previous test is deleted, if any
	Remove(key)

	created := Ingredient{Key: key, Name: name}

	//Save the new ingredient
	Save(&created)

	//Search it afterwards
	ingredient, err := FindByName(name)

	assert.NoError(t, err, "Load By Name error.")
	assert.NotEqual(t, ingredient, nil, "User is not found by name")

	assert.Equal(t, name, ingredient.Name, "Wrong name.")

	//purge the database with the test ingredient
	Remove(key)
}

func TestGetAll(t *testing.T) {
	ingredients, count, err := FindAll("name", "ASC", 0, 0)
	assert.NoError(t, err, "Error when loading ingredients")

	initialCount := len(ingredients)

	name1 := "TestGetAll1"
	name2 := "testGetAll2"
	key1 := "123"
	key2 := "456"
	ingredient1 := Ingredient{Name: name1, Key: key1}
	ingredient2 := Ingredient{Name: name2, Key: key2}

	Save(&ingredient1)
	Save(&ingredient2)

	ingredients, count, err = FindAll("name", "ASC", 0, 0)
	assert.Equal(t, len(ingredients), count, "Inconsistant data returned\n")
	assert.Equal(t, count, initialCount+2, "Ingredients not added in %v\n", ingredients)

	name1Found := false
	name2Found := false
	for _, ingredient := range ingredients {
		if ingredient.Name == name1 {
			if name1Found {
				assert.Fail(t, "Duplicate ingredient %v\n", name1)
			}
			name1Found = true
		}
		if ingredient.Name == name2 {
			if name2Found {
				assert.Fail(t, "Duplicate ingredient %v\n", name2)
			}
			name2Found = true
		}
	}
	assert.True(t, name1Found, "Wrong name for ingredient %v\n", name1)
	assert.True(t, name2Found, "Wrong name for ingredient %v\n", name2)

	assert.NoError(t, err, "Error when loading ingredients")

	Remove(key1)
	Remove(key2)

	ingredients, count, err = FindAll("name", "ASC", 0, 0)

	assert.Equal(t, len(ingredients), initialCount, "Ingredients not deleted from %v\n", ingredients)
	assert.NoError(t, err, "Ingredient deletion failed.")

}
