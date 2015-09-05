Feature:
    As a visitor of the website I want to see the Landing Page and be able to register my email

    Scenario: Displaying landing page@
  		When I GET /
  		Then response code should be 200
      Then response header Content-Type should be text/html
      Then response body should contain <input id="emailInput" type="email"

    Scenario: Deleting a contact
  		When I DELETE /contacts?email=testemail@mail.org
  		Then response code should be 204

    Scenario: Registering a contact
      Given I set body to email=testemail@mail.org
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
  		When I POST to /contacts
  		Then response code should be 201
      And response body should be valid json
      And response body path $.email should be testemail@mail.org
      And response body path $.created_at should be \w

    Scenario: Getting a contact
    	When I GET /contacts?email=testemail@mail.org
    	Then response code should be 200

    Scenario: Deleting a contact
  		When I DELETE /contacts?email=testemail@mail.org
  		Then response code should be 204

    Scenario: Getting a deleted contact
    	When I GET /contacts?email=testemail@mail.org
    	Then response code should be 404
