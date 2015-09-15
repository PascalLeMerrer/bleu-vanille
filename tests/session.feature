Feature:
    As a visitor of the website I can get connect and open a session, then disconnect

  Scenario: Sign up (creating an account) should succeed
    Given I set body to email=session_test1@mail.org;password=PASSWORD;firstname=JOHN;lastname=DOE
    And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
    When I POST to /users
    Then response code should be 201

  Scenario: I should no be able to get my profile without Authorization header
    When I GET /users/FAKE_USER_PROFILE/profile
    Then response code should be 401

  Scenario: I should no be able to get a user profile with an invalid Authorization header
    Given I set body to email=session_test1@mail.org;password=PASSWORD
    And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
    And I POST to /users/login
    And I store the value of header Content-Type as access token
    When I GET /users/FAKE_USER_PROFILE/profile
    Then response code should be 401

  Scenario: Sign in (login) should return an Authorization header
    Given I set body to email=session_test1@mail.org;password=PASSWORD
    And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
    When I POST to /users/login
    Then response code should be 200
    And response body should be valid json
    And response header Authorization should exist
    And I store the value of header Authorization as access token

    Scenario: Cleanup test data - Reconnect then delete account used for session test
      Given I set body to email=session_test1@mail.org;password=PASSWORD
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      And I POST to /users/login
      Then response code should be 200
      And I store the value of header Authorization as access token
      And I set bearer token
      When I POST to /users/delete
      Then response code should be 204

  # TODO Test logout then delete test account
