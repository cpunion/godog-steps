# Godog steps

[Godog](https://github.com/cucumber/godog) is a Cucumber implementation written in Go. It is a great tool for writing and running BDD tests.

This repository contains a set of steps that can be used in any project that uses Godog.

## Usage

To use these steps in your project, you need to import this repository as a Go module:

```bash
go get github.com/cpunion/godog-steps
```

Then, you can create youre `feature_test.go` file following the example at https://github.com/cucumber/godog#testmain and import the steps:

```go
package main

import (
	"os"
	"testing"

	steps "github.com/cpunion/godog-steps"
	"github.com/cpunion/godog-steps/file"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/cpunion/godog-steps/features"
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
		},
		Options: &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
```
