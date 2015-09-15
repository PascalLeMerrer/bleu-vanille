Feature:
  As an administrator, I shoul be allowed to perfomr admin tasks

  Scenario: Sign up (creating a basic account) should succeed
    Given I set body to email=admin_test1@mail.org;password=PASSWORD;firstname=JOHN;lastname=DOE
    And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
    When I POST to /users
    Then response code should be 201

  Scenario: I cannot not perform admin task with a non admin account
    Given I set body to email=admin_test1@mail.org;password=PASSWORD
    And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
    And I POST to /users/login
    Then response code should be 200
    And response header Authorization should exist
    And I store the value of header Authorization as access token
    And I set bearer token
  	When I GET /admin/contacts
  	Then response code should be 401

  Scenario: Authenticate as an admin
    When I set body to email=admin@bleuvanille.com;password=xeCuf8CHapreNe=
    And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
    When I POST to /users/login
    Then response code should be 200
    And response body should be valid json
    And response header Authorization should exist
    And I store the value of header Authorization as access token
    And I set bearer token

  Scenario: Cleanup test data - Reconnect then delete account used for admin test
    Given I set body to email=admin_test1@mail.org;password=PASSWORD
    And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
    And I POST to /users/login
    Then response code should be 200
    And I store the value of header Authorization as access token
    And I set bearer token
    When I POST to /users/delete
    Then response code should be 204
