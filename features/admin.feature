@admin
Feature:
  As an administrator, I shoul be allowed to perform admin tasks

  Scenario: Sign up (creating a basic account) should succeed
    Given I set body to email=admin_test1@mail.org;password=PASSWORD;firstname=JOHN_1;lastname=DOE_1
    And I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
    When I POST to /users
    Then response code should be 201
    And response body path $.id should be \w
    And I store the value of body path $.id as userId in global scope

  Scenario: Sign up (creating a basic account) should succeed
    Given I set body to email=admin_test2@mail.org;password=PASSWORD;firstname=JOHN_2;lastname=DOE_2
    And I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
    When I POST to /users
    Then response code should be 201
    And response body path $.id should be \w
    And I store the value of body path $.id as userId2 in global scope

  Scenario: Sign up (creating a basic account) should succeed
    Given I set body to email=admin_test3@mail.org;password=PASSWORD;firstname=JOHN_3;lastname=DOE_3
    And I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
    When I POST to /users
    Then response code should be 201
    And response body path $.id should be \w
    And I store the value of body path $.id as userId3 in global scope

  Scenario: Sign up (creating a basic account) should succeed
    Given I set body to email=admin_test4@mail.org;password=PASSWORD;firstname=JOHN_4;lastname=DOE_4
    And I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
    When I POST to /users
    Then response code should be 201
    And response body path $.id should be \w
    And I store the value of body path $.id as userId4 in global scope


  Scenario: I cannot not perform admin task with a non admin account
    Given I set body to email=admin_test1@mail.org;password=PASSWORD
    And I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
    When I POST to /users/login
    Then response code should be 200
    And response header Authorization should exist
    And I store the value of header Authorization as access token
    And I set bearer token
    When I GET /admin/contacts
    Then response code should be 401


  Scenario: Authenticate as an admin
    Given I set Content-Type header to application/x-www-form-urlencoded;charset=UTF-8
    When I POST to /users/login
    Then response code should be 200
    And response body should be valid json
    And response header Authorization should exist
    And I store the value of header Authorization as access token
    And I set bearer token


  Scenario: As an admin, Get all users
    Given I set bearer token
    And I set Content-Type header to application/json;charset=UTF-8
    And I set Accept header to application/json
    When I GET /admin/users
    Then response code should be 200
    And response header X-TOTAL-COUNT should exist
    And response header X-TOTAL-COUNT should be \d+
    And response body should contain "email": "admin@bleuvanille.com"
    And response body should be valid json
    And response body path $[0].id should be \w
    And response body path $[0].firstname should be \w
    And response body path $[0].lastname should be \w
    And response body path $[0].email should be \w
    And response body path $[0].createdAt should be \w

  Scenario: As an admin, get a limited subset of user list
    Given I set bearer token
    And I set Content-Type header to application/json;charset=UTF-8
    And I set Accept header to application/json
    When I GET /admin/users?offset=1&limit=2
    Then response code should be 200
    And response header X-TOTAL-COUNT should exist
    And response header X-TOTAL-COUNT should be \d+
    And response body should be valid json
    And response body at path $.* should be a json array
    And response body at path $.* should be an array of length 2
    And response body path $[0].id should be \w
    And response body path $[0].firstname should be \w
    And response body path $[0].lastname should be \w
    And response body path $[0].email should be \w
    And response body path $[0].createdAt should be \w
    And response body path $[1].id should be \w
    And response body path $[1].firstname should be \w
    And response body path $[1].lastname should be \w
    And response body path $[1].email should be \w
    And response body path $[1].createdAt should be \w

  Scenario: As an admin, search for a user with a given email
    Given I set bearer token
    And I set Content-Type header to application/json;charset=UTF-8
    And I set Accept header to application/json
    When I GET /admin/users/email?email=admin_test1@mail.org
    Then response code should be 200
    And response body should be valid json
    And response body path $.id should be \w
    And response body path $.firstname should be JOHN
    And response body path $.lastname should be DOE
    And response body path $.email should be admin_test1@mail.org
    And response body path $.createdAt should be \w


  Scenario: check rights on account used for admin tests
    Given I set bearer token
    And I set Content-Type header to application/json;charset=UTF-8
    When I GET /users/`userId`
    Then response code should be 200
    And response body path $.isAdmin should be false


  Scenario: modifying email on account should not modify rights
    Given I set bearer token
    And I set Content-Type header to application/json;charset=UTF-8
    When I PATCH /users/`userId` with body
    """
    {
        "email": "admin_test5@mail.org",
        "firstname": "Alphonse",
        "lastname": "Dans l'tas"
    }
    """
    Then response code should be 200
    When I GET /users/`userId`
    Then response code should be 200
    And response body path $.isAdmin should be false
    And response body path $.email should be admin_test5@mail.org
    And response body path $.firstname should be Alphonse
    And response body path $.lastname should be Dans l'tas


  Scenario: modifying rights on account used for admin tests
    Given I set bearer token
    And I set Content-Type header to application/json;charset=UTF-8
    When I PATCH /users/`userId` with body
    """
    {
        "isAdmin": true
    }
    """
    Then response code should be 200
    When I GET /users/`userId`
    Then response code should be 200
    And response body path $.isAdmin should be true
    When I PATCH /users/`userId` with body
    """
    {
        "isAdmin": false
    }
    """
    Then response code should be 200
    When I GET /users/`userId`
    Then response code should be 200
    And response body path $.isAdmin should be false


  Scenario: Cleanup test data - delete account used for admin test
    Given I set bearer token
    And I set Content-Type header to application/json;charset=UTF-8
    When I DELETE /admin/users/`userId`
    Then response code should be 204

  Scenario: Cleanup test data - delete account used for admin test
    Given I set bearer token
    And I set Content-Type header to application/json;charset=UTF-8
    When I DELETE /admin/users/`userId2`
    Then response code should be 204

  Scenario: Cleanup test data - delete account used for admin test
    Given I set bearer token
    And I set Content-Type header to application/json;charset=UTF-8
    When I DELETE /admin/users/`userId3`
    Then response code should be 204

  Scenario: Cleanup test data - delete account used for admin test
    Given I set bearer token
    And I set Content-Type header to application/json;charset=UTF-8
    When I DELETE /admin/users/`userId4`
    Then response code should be 204
