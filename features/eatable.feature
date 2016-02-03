@eatable
Feature:  
    As a admin of the website, I want to manage my eatables

Scenario: Creating an eatable without user should fail 
    When I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    And I POST to /eatables  with body 
    """
    {
        "name" : "pomme de terre",
        "type" : "ingredient",
        "description" : "La pomme de terre est un légume découvert aux amériques."
    }
    """
    Then response code should be 401 
    
Scenario: Creating an eatable with authenticated user 
    Given I log as test user 
    And I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I POST to /eatables with body 
    """
    {
        "name" : "pomme de terre",
        "type" : "ingredient",
        "description" : "La pomme de terre est un légume découvert aux amériques.",
        "status" : "statusinterditamettreajouralacreation"
    }
    """
    Then response code should be 201     
    And response body should be valid json 
    And response body path $.name should be pomme de terre 
    And response body path $.type should be ingredient 
    And response body path $.status should be new
    And response body path $.description should be La pomme de terre est un légume découvert aux amériques. 
    And response body path $.createdAt should be \w 
    And response body path $._id should be \w 
    And response body path $._key should be \d+ 
    And I store the value of body path $._key as eatableKey in global scope 
    
Scenario: Creating an eatable with authenticated user 
    Given I log as test user 
    And I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I POST to /eatables with body 
    """
    {
        "name" : "pomme de terre",
        "type" : "ingredient",
        "description" : "La pomme de terre est un légume découvert aux amériques.",
        "status" : "statusinterditamettreajouralacreation"
    }
    """
    Then response code should be 409
    And the JSON should be
    """
    { 
        "error" : "Eatable with name pomme de terre already exists."
    }
    """ 

Scenario: Modifying an eatable with user with a wrong type
    Given I set Cookie header to global variable cookie 
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json
    When I PUT /eatables/`eatableKey` with body 
    """
    {
    "name" : "pomme de terre",
    "type" : "ingredientinconnu",
        "status" : "active",
        "description" : "La pomme de terre est un légume découvert aux amériques."
    }
    """
    Then response code should be 400 
    And response body should be valid json 
    And response body path $.error should be Invalid eatable type or name 

Scenario: Modifying an eatable with user 
    Given I set Cookie header to global variable cookie
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json 
    When I PUT /eatables/`eatableKey` with body 
    """
    {
        "name" : "pomme de terre nouvelle",
        "type" : "ingredientrecette",
        "description" : "La pomme de terre est un légume découvert aux amériques. Elle a sauvé la vie à des millions d Irlandais."
    }
    """
    Then response code should be 200
    When I GET /eatables/`eatableKey`
    Then response code should be 200
    And response body should be valid json
    And response body path $.name should be pomme de terre nouvelle
    And response body path $.type should be ingredient 
    And response body path $.status should be new
    And response body path $.description should be La pomme de terre est un légume découvert aux amériques. Elle a sauvé la vie à des millions d Irlandais.
    And response body path $.createdAt should be \w 
    And response body path $._id should be \w    
    And response body path $._key should be \d+    
        
Scenario: Setting the nutrient of an eatable with admin user 
    Given I set Cookie header to global variable cookie
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json 
    When I PUT /eatables/`eatableKey`/nutrient with body
    """
    {
        "carbohydrate" : 10,
        "sugar" : 9,
        "protein" : 11,
        "lipid" : 12,
        "fiber" : 1,
        "alcohol" : 2
    }
    """
    Then response code should be 200 
    
Scenario: Changing the parent of an eatable with admin user from its id 

    Given I set Cookie header to global variable cookie
    And I set Content-Type header to application/json; charset=UTF-8
    And I set Accept header to application/json 
    When I POST to /eatables with body 
    """
    {
        "name" : "légume",
        "type" : "ingredient"
    }
    """
    Then  response code should be 201
    And   I store the value of body path $._key as parentKey in global scope
    Given I set Cookie header to global variable cookie  
    When  I PUT /eatables/`eatableKey`/parent/`parentKey`
    Then  response code should be 200 
    And   response body should be valid json 

Scenario: Get the full eatable data 
    Given I set Cookie header to global variable cookie
    And   I set Content-Type header to application/json; charset=UTF-8
    And   I set Accept header to application/json 
    When  I GET /eatables/`eatableKey`
    Then  response code should be 200 
    And   response body should be valid json 
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

Scenario: Get eatable with an unkown id
   Given I set Cookie header to global variable cookie 
   When  I GET /eatables/1
   Then  response code should be 404 
   And   response body should be valid json
   And   the JSON should be
   """
   {"error":"No eatable found for key: 1"}
   """
    
Scenario: Disabling an eatable with admin user 
  Given I log as admin user
  And I set Cookie header to global variable cookie 
  When I PATCH /admin/eatables/`eatableKey`/status with body
  """
  { "status": "disabled" }
  """  
  Then response code should be 200 
  And response body should be valid json 
  And response body path $.status should be disabled 
  Given I PATCH /admin/eatables/`eatableKey`/status with body
  """
  { "status": "active" }
  """ 
  Then response code should be 200 
  And response body should be valid json 
  And response body path $.status should be active 

    
Scenario: Deleting an eatable with admin user 
  Given I log as admin user
  And I set Cookie header to global variable cookie 
  When I DELETE /admin/eatables/`eatableKey`
  Then response code should be 204 
  When I GET /eatables/`eatableKey`
  Then response code should be 404
  When I DELETE /admin/eatables/`parentKey`
  Then response code should be 204 
  When I GET /eatables/`parentKey`
  Then response code should be 404
