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
import (
  // ... other code
	steps "github.com/cpunion/godog-steps"
	"github.com/cpunion/godog-steps/file"
	"github.com/cucumber/godog"
	// ... other code
)

// ... other code

func TestMain(m *testing.M) {
  // ... other code
	file.InitAssetsFS(features.AssetsFS) // Unnessesary if you don't use files mock

	status := godog.TestSuite{
		Name: "godogs",
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			steps.Init(s) // Initialize steps in this REPO
      // ... Your steps
		},
		Options: &opts,
	}.Run()
  // ... other code
}
```
