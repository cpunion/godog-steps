Feature: run godog-steps demo

  Scenario: test file godog-steps
    Given I have a file "a.txt" with content:
      """
      Hello world!
      Hello godog!
      """
    When I run "cat a.txt"
    Then I should see output "Hello world!"
    And I should see output "Hello godog!"
    And I am happy
