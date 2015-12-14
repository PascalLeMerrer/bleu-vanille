@eatable
Feature:  
	As a admin of the website, I want to manage my eatables

Scenario: Creating an eatable without user should fail 
	When I POST to /eatables  with body 
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
	And response body path $.created_at should be \w 
	And response body path $._id should be \w 
	And I store the value of body path $._key as createdIngredientId in global scope 
	
Scenario: Modifying an eatable with user with a wrong type
	Given I set Cookie header to global variable cookie 
	When I PUT /eatables/`createdIngredientId` with body 
		"""
	  		{
				"name" : "pomme de terre",
				"type" : "ingredientinconnu",
	  			"status" : "active",
	  			"description" : "La pomme de terre est un légume découvert aux amériques. Elle a sauvé la vie à des millions d Irlandais."
	  		}
	  		"""
	Then response code should be 400 
	And response body should be valid json 
	And response body path $.error should be Unknown type : ingredientinconnu 

Scenario: Modifying an eatable with user 
	Given I set Cookie header to global variable cookie 
	When I PUT /eatables/`createdIngredientId` with body 
		"""
	  		{
	  			"name" : "pomme de terre nouvelle",
	  			"type" : "ingredientrecette",
	  			"description" : "La pomme de terre est un légume découvert aux amériques. Elle a sauvé la vie à des millions d Irlandais."
	  		}
	  		"""
	Then response code should be 200 
		
Scenario: Setting the nutrient of an eatable with admin user 
	Given I set Cookie header to global variable cookie 
	When I PUT /eatables/`createdIngredientId`/nutrient with body
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
	
#Scenario: Disabling an eatable with admin user 
#	Given I log as admin user 
#	When I PATCH /eatables/`createdIngredientId`/status/disabled 
#	Then  response code should be 200 
#	And  response body should be valid json 
#	And response body path $.status should be ok 
#	Given I log as admin user 
#	When I PATCH /eatables/`createdIngredientId`/status/active 
#	Then  response code should be 200 
#	And  response body should be valid json 
#	And response body path $.status should be ok 
	
Scenario: Changing the parent of an eatable with admin user from its id 
	Given I set Cookie header to global variable cookie 
	When I POST to /eatables with body 
		"""
	  		{
				"name" : "légume",
				"type" : "ingredient"
	  		}
	  		"""
	Then response code should be 201
    And I store the value of body path $._key as createdIngredientParentId in global scope
    And I store the value of body path $._key as createdIngredientParentGlobalId in global scope
	Given I set Cookie header to global variable cookie  
	When  I PUT /eatables/`createdIngredientId`/parent/`createdIngredientParentId`
	Then  response code should be 200 
	And  response body should be valid json 

Scenario: Get the full eatable data 
	Given I set Cookie header to global variable cookie 
	When   I GET /eatables/`createdIngredientId`
	Then   response code should be 200 
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
		  			}
	  			}
	"""
   # Desactivé car response body ne parse par la variable (même globale)
   #And response body path $.parentid should be eatables/`createdIngredientParentGlobalId`

#Scenario: Get eatable with an unkown id
#    Given I set Cookie header to global variable cookie 
#    When   I GET /eatables/1
#    Then   response code should be 404 
#    And response body path $.error should be Unknown Id : 1