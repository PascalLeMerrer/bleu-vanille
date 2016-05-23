@ingredient
Feature:
    As a user of the website I can define ingredients and get info about them

    Scenario: Registering an ingredient
      Given I log as test user
      And I store the value of body path $.id as userId in scenario scope
      And I set Content-Type header to application/json
      When I POST to /ingredients with body
      """
      {
        "name": "un nom",
        "energy": 123,
        "category": "viande",
        "months": [5,6,7]
      }
      """
      Then response code should be 201
      And response header Content-Type should be application/json
      And response body should be valid json
      And response body path $.createdAt should be \w
      And response body path $._key should be \d+
      And I store the value of body path $._key as ingredientId in global scope
      And the JSON should be
      """
      {
        "name": "un nom",
        "energy": 123,
        "category": "viande",
        "months": [5,6,7],
        "approved" : false,
        "creator": "`userId`"
      }
      """

    Scenario: Get a given ingredient
      Given I send the cookie token
      When I GET /ingredients/`ingredientId`
      Then response code should be 200
      And response header Content-Type should be application/json
      And response body should be valid json
      And response body path $._key should be \d+
      And response body path $.createdAt should be \w
      And the JSON should be
      """
      {
        "name": "un nom",
        "energy": 123,
        "category": "viande",
        "months": [5,6,7],
        "approved" : false
      }
      """

    Scenario: Registering the same ingredient twice should return a Conflict HTTP Error
      Given I send the cookie token
      And I set Content-Type header to application/json
      When I POST to /ingredients with body
      """
      {
        "name": "un nom",
        "energy": 123,
        "category": "viande",
        "months": [5,6,7]
      }
      """
      Then response code should be 409

    Scenario: A user can modify an ingredient he created, before it was approved
      Given I send the cookie token
      And I set Content-Type header to application/json
  		When I PATCH /ingredients/`ingredientId` with body
      """
      {
        "energy":456,
        "category":"légume",
        "months":[7,8,9]
      }
      """
  		Then response code should be 200
      And the JSON should be
      """
      {
        "name": "un nom",
        "energy": 456,
        "category": "légume",
        "months": [7,8,9],
        "approved": false
      }
      """

    Scenario: A simple user cannot approve an ingredient
      Given I send the cookie token
      And I set Content-Type header to application/json
      When I PATCH /ingredients/`ingredientId` with body
      """
      {
        "approved": true
      }
      """
      Then response code should be 403
      When I GET /ingredients/`ingredientId`
      Then response code should be 200
      And the JSON should be
      """
      {
        "name": "un nom",
        "energy": 456,
        "category": "légume",
        "months": [7,8,9],
        "approved": false
      }
      """

    Scenario: An admin can approve an ingredient
      Given I log as admin user
      And I set Content-Type header to application/json
      When I PATCH /ingredients/`ingredientId` with body
      """
      {
        "approved": true
      }
      """
      Then response code should be 200
      And the JSON should be
      """
      {
        "name": "un nom",
        "energy": 456,
        "category": "légume",
        "months": [7,8,9],
        "approved": true
      }
      """

    Scenario: As an admin, I cannot delete an approved ingredient
      Given I send the cookie token
      And I DELETE /ingredients/`ingredientId`
      Then response code should be 403

    Scenario: A simple user cannot unapprove an ingredient
      Given I log as test user
      And I set Content-Type header to application/json
      When I PATCH /ingredients/`ingredientId` with body
      """
      {
        "approved": false
      }
      """
      Then response code should be 403
      When I GET /ingredients/`ingredientId`
      Then response code should be 200
      And the JSON should be
      """
      {
        "name": "un nom",
        "energy": 456,
        "category": "légume",
        "months": [7,8,9],
        "approved": true
      }
      """

    Scenario: An admin can unapprove an ingredient
      Given I log as admin user
      And I set Content-Type header to application/json
      When I PATCH /ingredients/`ingredientId` with body
      """
      {
        "approved": false
      }
      """
      Then response code should be 200
      And the JSON should be
      """
      {
        "name": "un nom",
        "energy": 456,
        "category": "légume",
        "months": [7,8,9],
        "approved": false
      }
      """

    Scenario: As a simple user, I cannot delete an unapproved ingredient I did not create
      Given I send the cookie token
      And I set Content-Type header to application/json
      When I POST to /ingredients with body
      """
      {
        "name": "oignon",
        "energy": 45,
        "category": "légume",
        "months": [5,6,7,8,9,10]
      }
      """
      Then response code should be 201
      And response header Content-Type should be application/json
      And response body should be valid json
      And response body path $._key should be \d+
      And I store the value of body path $._key as onionId in global scope
      Given I clear the cookie token
      And I log as test user
      And I DELETE /ingredients/`onionId`
      Then response code should be 403


    Scenario: As an admin, Get all ingredients
      Given I log as admin user
      And I set Content-Type header to application/json;charset=UTF-8
      And I set Accept header to application/json
      When I GET /ingredients
      Then response code should be 200
      And response body should contain "name": "un nom"
      And response body should contain "name": "oignon"
      And response header Content-Type should be application/json


    Scenario: As an admin, I can delete an unapproved ingredient
      Given I send the cookie token
      And I DELETE /ingredients/`ingredientId`
      Then response code should be 204
      And I set Accept header to application/json
      And I GET /ingredients/`ingredientId`
      Then response code should be 404
      And I DELETE /ingredients/`onionId`
      Then response code should be 204
      And I set Accept header to application/json
      And I GET /ingredients/`onionId`
      Then response code should be 404


    Scenario: As an admin, deleting a non existing ingredient should return an HTTP code 404
      Given I send the cookie token
      And I DELETE /ingredients/99999999
      Then response code should be 404