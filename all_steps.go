package steps

import (
	"github.com/cpunion/godog-steps/browser"
	"github.com/cpunion/godog-steps/cmd"
	"github.com/cpunion/godog-steps/common"
	"github.com/cpunion/godog-steps/file"
	"github.com/cpunion/godog-steps/image"
	"github.com/cpunion/godog-steps/network"
	"github.com/cucumber/godog"
)

func Init(ctx *godog.ScenarioContext) {
	cmd.Init(ctx)
	file.Init(ctx)
	browser.Init(ctx)
	image.Init(ctx)
	network.Init(ctx)
	common.Init(ctx)
}
