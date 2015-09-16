Feature:
    As a visitor of the website I create an account, connect, disconnect and manage my account

    Scenario: Sign up (creating an account)
      Given I set body to email=user_test1@mail.org;password=mypassword;firstname=John;lastname=Doe
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
  		When I POST to /users
  		Then response code should be 201
      And response body should be valid json
      And response body path $.firstname should be John
      And response body path $.lastname should be Doe
      And response body path $.email should be user_test1@mail.org
      And response body path $.createdAt should be \w

    Scenario: Sign up with already registered email should result in conflict
      Given I set body to email=user_test1@mail.org;password=PASSWORD;firstname=JOHN;lastname=DOE
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      When I POST to /users
      Then response code should be 409

    Scenario: Sign in (login)
      Given I set body to email=user_test1@mail.org;password=mypassword
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      When I POST to /users/login
      Then response code should be 200
      And response body should be valid json
      And response header Authorization should exist
      And response body path $.token should be \w

    Scenario: After I signed in, I should be able to get my profile
      Given I set body to email=user_test1@mail.org;password=mypassword
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      And I POST to /users/login
      And I store the value of header Authorization as access token
      And I set bearer token
      And I store the value of body path $.id as userId in scenario scope
      When I GET /users/`userId`/profile
      Then response code should be 200
      # TODO: to complete


    Scenario: Sign in with wrong email should fail
      Given I set body to email=fake@mail.org;password=mypassword
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      When I POST to /users/login
      Then response code should be 401
      And response header Set-Cookie should not exist
      And response header Authorization should not exist
      And response body should not contain token

    Scenario: Sign in with wrong password should fail
      Given I set body to email=user_test1@mail.org;password=BAD_PASSWORD
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      When I POST to /users/login
      Then response code should be 401
      And response header Set-Cookie should not exist
      And response header Authorization should not exist
      And response body should not contain token

    Scenario: Delete account with wrong password should fail
      Given I set body to email=user_test1@mail.org;password=BAD_PASSWORD
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      When I POST to /users/delete
      Then response code should be 401
      And response header Set-Cookie should not exist
      And response header Authorization should not exist
      And response body should not contain token

    Scenario: Delete account
      Given I set body to email=user_test1@mail.org;password=mypassword
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      And I POST to /users/login
      And I store the value of header Authorization as access token
      And I set bearer token
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      And I set body to password=mypassword
      When I POST to /users/delete
      Then response code should be 204
      And response header Set-Cookie should not exist
      And response header Authorization should not exist
      And response body should not contain token

    Scenario: Sign in with deleted account credentials should fail
      Given I set body to email=user_test1@mail.org;password=mypassword
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      When I POST to /users/login
      Then response code should be 401
      And response header Set-Cookie should not exist
      And response header Authorization should not exist
      And response body should not contain token

    Scenario: Create an account for password change
      Given I set body to email=user_test2@mail.org;password=OLDPASSWORD;firstname=John;lastname=Doe
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
  		And I POST to /users
  		Then response code should be 201

    Scenario: Sign in with initial password
      Given I set body to email=user_test2@mail.org;password=OLDPASSWORD
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      When I POST to /users/login
      Then response code should be 200
      And response header Authorization should exist
      And I store the value of header Authorization as access token

    Scenario: Change password
      Given I set body to email=user_test2@mail.org;oldPassword=OLDPASSWORD;newPassword=NEWPASSWORD
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      And I set bearer token
  		And I PUT /users/password
      Then response code should be 200

    Scenario: Sign in with new password (login)
      Given I set body to email=user_test2@mail.org;password=NEWPASSWORD
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      When I POST to /users/login
      Then response code should be 200
      And response body should be valid json
      And response header Set-Cookie should not exist
      And response body path $.token should be \w
      And response header Authorization should exist
      And I store the value of header Authorization as access token

    Scenario: I should not be allowed to delete my account with a wrong password
      Given I set body to password=BADPASSWORD
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      And I set bearer token
      When I POST to /users/delete
      Then response code should be 401

    Scenario: Delete the account used for password change
      Given I set body to password=NEWPASSWORD
      And I set Content-Type header to application/x-www-form-urlencoded; charset=UTF-8
      And I set bearer token
      When I POST to /users/delete
      Then response code should be 204
      And response header Set-Cookie should not exist
      And response body should not contain token
