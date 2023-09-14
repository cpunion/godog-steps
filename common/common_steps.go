package common

import (
	"context"
	"time"

	"github.com/cucumber/godog"
)

// for debugging
func iWaitSeconds(ctx context.Context, seconds int) (context.Context, error) {
	time.Sleep(time.Duration(seconds) * time.Second)
	return ctx, nil
}

func Init(ctx *godog.ScenarioContext) {
	ctx.Step(`^I wait (\d+) seconds?$`, iWaitSeconds)
}
