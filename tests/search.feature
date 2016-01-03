@search
Feature:  
	As an user of the website, I want to search for eatables
	
Scenario: Creating and searching for an eatable with authenticated user 
    Given I log as test user 
    And I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I POST to /eatables with body 
    """
    {
        "name" : "carotte",
        "type" : "ingredient",
        "description" : "La carotte aime bien le lapin"
    }
    """
    Then response code should be 201
    And response body path $._key should be \d+
    And I store the value of body path $._key as eatableKey in global scope
    Given I set Cookie header to global variable cookie
    When I GET /search/index/`eatableKey`
    Then response code should be 200
    Given I set Cookie header to global variable cookie 
    When I GET /search/carotte
    Then response code should be 200
    And response body should be valid json
    And   the JSON should be
    """
    {
        "name" : "pomme de terre nouvelle",
        "type" : "ingredientrecette",
        "status" : "new",
        "description" : "La pomme de terre est un légume découvert aux amériques. Elle a sauvé la vie à des millions d Irlandais.",
        "nutrient" : {
            "carbohydrate" : 10,
            "sugar" : 9,
            "protein" : 11,
            "lipid" : 12,
            "fiber" : 1,
            "alcohol" : 2
        },
        "parent": {
            "name": "légume",
            "status": "new",
            "type": "ingredient"
        }
    }
    """
    