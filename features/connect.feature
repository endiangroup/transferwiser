Feature: Connect the service with a Transferwise account
  As a business owner
  I want to be able to connect the service with my Transferwise account

  Scenario:	Transferwise connect link
    Given the service has not been authenticated with Transferwise
    When I visit the link to connect with Transferwise
    Then I should be redirected to the Transferwise authorization login page

  Scenario: Return from Transferwise connect link
    Given Transferwise refresh token response for 'abcd1234' is:
      | access_token  | myaccesstoken |
      | token_type    | bearer        |
      | refresh_token | abcd1234      |
      | expires_in    | 3600          |
      | scope         | transfers     |
    And the service has not been authenticated with Transferwise
    When I return from Transferwise OAuth with code 'abcd1234'
    Then the service is authenticated with Transferwise
