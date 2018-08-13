Feature: Connect the service with a Transferwise account
  As a business owner
  I want to be able to connect the service with my Transferwise account

  Scenario:	Oauth connect link
    Given The service has not been authenticated with Transferwise
    When I visit the link to connect with Transferwise
    Then I should be redirected to the Transferwise Oauth login page
