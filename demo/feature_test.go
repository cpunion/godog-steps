package main

import (
	"os"
	"testing"

	steps "github.com/cpunion/godog-steps"
	"github.com/cpunion/godog-steps/demo/features"
	"github.com/cpunion/godog-steps/demo/features/step_definitions"
	"github.com/cpunion/godog-steps/file"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts) // godog v0.11.0 and later
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	file.InitAssetsFS(features.AssetsFS)

	status := godog.TestSuite{
		Name: "godogs",
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			steps.Init(s)
			step_definitions.InitializeScenarioContext(s)
		},
		Options: &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
