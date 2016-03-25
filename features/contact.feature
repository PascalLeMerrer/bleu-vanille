@contact
Feature:
    As a visitor of the website I want to see the Landing Page and be able to register my email

    Scenario: Displaying landing page
  		When I GET /
  		Then response code should be 200
      Then response header Content-Type should be text/html
      Then response body should contain input id="emailInput" type="email"

    Scenario: Registering a contact
      Given I set body to email=testemail@gmail.com
      And I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
      When I POST to /contacts
      Then response code should be 201
      And response body should be valid json
      And response body path $.email should be testemail@gmail.com
      And response body path $.created_at should be \w

    Scenario: Registering a contact with an email address including a tag
      # name+tag@gmail.com must be url encoded in the post body
      Given I set body to email=name%2Btag@gmail.com
      And I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
      When I POST to /contacts
      Then response code should be 201
      And response body should be valid json
      And response body path $.email should be name\+tag@gmail.com
      And response body path $.created_at should be \w

    Scenario: Registering a contact with a bad email address should fail
      Given I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
      And I set body to email=me@
      When I POST to /contacts
      Then response code should be 400

    Scenario: Registering a contact with a bad email address should fail
      Given I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
      And I set body to email=@gmail.com
      When I POST to /contacts
      Then response code should be 400

    Scenario: Registering a contact with a bad email address should fail
      Given I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
      And I set body to email=me@example..com
      When I POST to /contacts
      Then response code should be 400

    Scenario: Registering a contact with a bad email address should fail
      Given I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
      And I set body to email=me.example@com
  		When I POST to /contacts
  		Then response code should be 400

    Scenario: Registering the same contact twice should returns a Conflict HTTP Error
      Given I set body to email=testemail@gmail.com
      And I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
  		When I POST to /contacts
  		Then response code should be 409

    Scenario: Authenticate as an admin
      When I log as admin user
      And I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
      When I POST to /users/login
      Then response code should be 200
      And response body should be valid json
      And response header Authorization should exist
      And I store the value of header Authorization as access token
      And I set bearer token

    Scenario: As an admin, Verify contact is registered
      When I set bearer token
      And I set Content-Type header to application/json;charset=UTF-8
      And I set Accept header to application/json
      And I GET /admin/contacts
      Then response code should be 200
      And response body should contain testemail@gmail.com

    Scenario: As an admin, Get all contacts
      When I set bearer token
      And I set Content-Type header to application/json;charset=UTF-8
      And I set Accept header to application/json
      And I GET /admin/contacts
      Then response code should be 200
      And response body should contain testemail@gmail.com
      And response header Content-Type should be application/json

    Scenario: As an admin, Download contacts
      When I set bearer token
      And I set Accept header to text/csv
      And I GET /admin/contacts
      Then response code should be 200

    Scenario: As an admin, Deleting a non existing contact should return an HTTP code 204
      Given I set bearer token
  		And I DELETE /admin/contacts?email=unknown@gmail.com
  		Then response code should be 204

    Scenario: As an admin, Deleting a contact
      Given I set bearer token
  		And I DELETE /admin/contacts?email=testemail@gmail.com
  		Then response code should be 204
      And I set Accept header to application/json
      And I GET /admin/contacts
    	Then response code should be 200
      And response body should not contain testemail@gmail.com

    Scenario: As an admin, Deleting test contact
      Given I set bearer token
      And I DELETE /admin/contacts?email=name%2Btag@gmail.com
      Then response code should be 204


