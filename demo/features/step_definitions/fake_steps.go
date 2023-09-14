package step_definitions

import "github.com/cucumber/godog"

func iAmHappy() error {
	return nil
}

func InitializeScenarioContext(ctx *godog.ScenarioContext) {
	ctx.Step(`^I am happy$`, iAmHappy)
}
