@admin
Feature:
  As an administrator, I shoul be allowed to perform admin tasks

  Scenario: Sign up (creating a basic account) should succeed
    Given I set body to email=admin_test1@mail.org;password=PASSWORD;firstname=JOHN;lastname=DOE
    And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
    When I POST to /users
    Then response code should be 201
    And response body path $.id should be \w
    And I store the value of body path $.id as userId in global scope


  Scenario: I cannot not perform admin task with a non admin account
    Given I set body to email=admin_test1@mail.org;password=PASSWORD
    And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
    When I POST to /users/login
    Then response code should be 200
    And response header Authorization should exist
    And I store the value of header Authorization as access token
    And I set bearer token
    When I GET /admin/contacts
    Then response code should be 401

  Scenario: Authenticate as an admin
    Given I set body to email=admin@bleuvanille.com;password=xeCuf8CHapreNe=
    And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
    When I POST to /users/login
    Then response code should be 200
    And response body should be valid json
    And response header Authorization should exist
    And I store the value of header Authorization as access token
    And I set bearer token

  Scenario: As an admin, Get all users
    Given I set bearer token
    And I set Content-Type header to application/json; charset=UTF-8
    When I GET /admin/users
    Then response code should be 200
    And response body should contain "email":"admin@bleuvanille.com"
    And response body should be valid json
    And response body path $[0].id should be \d+
    And response body path $[0].firstname should be \w
    And response body path $[0].lastname should be \w
    And response body path $[0].email should be \w
    And response body path $[0].createdAt should be \w


  Scenario: Cleanup test data - delete account used for admin test
    Given I set bearer token
    And I set Content-Type header to application/json; charset=UTF-8
    When I DELETE /admin/users/`userId`
    Then response code should be 204
