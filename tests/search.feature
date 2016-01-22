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
    [
    	{
        	"name" : "carotte"
    	}
    ]
    """

Scenario: Deleting eatables created for the test
  Given I log as admin user
  And I set Cookie header to global variable cookie 
  When I DELETE /admin/unindex/`eatableKey`
  Then response code should be 204
  Given I log as admin user
  And I set Cookie header to global variable cookie 
  When I DELETE /admin/eatables/`eatableKey`
  Then response code should be 204
  
Scenario: Creating and searching for many eatables with authenticated user : Legume1, Légume2, Choux-fleur, Asperge, Asperges, Asperge2, Asperge3, Asperge4
    Given I log as test user 
    And I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I POST to /eatables with body 
    """
    {
        "name" : "Légume1",
        "type" : "ingredient",
        "description" : ""
    }
    """
    Then response code should be 201
    And response body path $._key should be \d+
    And I store the value of body path $._key as eatableKey in global scope
    Given I set Cookie header to global variable cookie
    When I GET /search/index/`eatableKey`
    Then response code should be 200
    
    Given I log as test user 
    And I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I POST to /eatables with body 
    """
    {
        "name" : "Choux-fleur",
        "type" : "ingredient",
        "description" : ""
    }
    """
    Then response code should be 201
    And response body path $._key should be \d+
    And I store the value of body path $._key as eatableKey in global scope
    Given I set Cookie header to global variable cookie
    When I GET /search/index/`eatableKey`
    Then response code should be 200
    
    Given I log as test user 
    And I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I POST to /eatables with body 
    """
    {
        "name" : "Asperge",
        "type" : "ingredient",
        "description" : ""
    }
    """
    Then response code should be 201
    And response body path $._key should be \d+
    And I store the value of body path $._key as eatableKey in global scope
    Given I set Cookie header to global variable cookie
    When I GET /search/index/`eatableKey`
    Then response code should be 200

    Given I log as test user 
    And I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I POST to /eatables with body 
    """
    {
        "name" : "Asperge2",
        "type" : "ingredient",
        "description" : ""
    }
    """
    Then response code should be 201
    And response body path $._key should be \d+
    And I store the value of body path $._key as eatableKey in global scope
    Given I set Cookie header to global variable cookie
    When I GET /search/index/`eatableKey`
    Then response code should be 200

    Given I log as test user 
    And I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I POST to /eatables with body 
    """
    {
        "name" : "Asperge3",
        "type" : "ingredient",
        "description" : ""
    }
    """
    Then response code should be 201
    And response body path $._key should be \d+
    And I store the value of body path $._key as eatableKey in global scope
    Given I set Cookie header to global variable cookie
    When I GET /search/index/`eatableKey`
    Then response code should be 200

    Given I log as test user 
    And I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I POST to /eatables with body 
    """
    {
        "name" : "Asperge4",
        "type" : "ingredient",
        "description" : ""
    }
    """
    Then response code should be 201
    And response body path $._key should be \d+
    And I store the value of body path $._key as eatableKey in global scope
    Given I set Cookie header to global variable cookie
    When I GET /search/index/`eatableKey`
    Then response code should be 200
    
    Given I log as test user 
    And I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I POST to /eatables with body 
    """
    {
        "name" : "Légume2",
        "type" : "ingredient",
        "description" : ""
    }
    """
    Then response code should be 201
    And response body path $._key should be \d+
    And I store the value of body path $._key as eatableKey in global scope
    Given I set Cookie header to global variable cookie
    When I GET /search/index/`eatableKey`
    Then response code should be 200
    
    Given I set Cookie header to global variable cookie 
    When I GET /search/all
    Then response code should be 200
    And response body should be valid json
	And response header X-TOTAL-COUNT should be 7