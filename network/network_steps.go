package network

import (
	"context"
	"net"
	"time"

	"github.com/cpunion/godog-steps/helper"
	"github.com/cucumber/godog"
)

func itShouldBeRunningOnPort(ctx context.Context, addr string) (context.Context, error) {
	toCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	helper.WaitFor(toCtx, func() bool {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	})
	return ctx, nil
}

func Init(ctx *godog.ScenarioContext) {
	ctx.Step(`^The server should be running on "([^"]*)"$`, itShouldBeRunningOnPort)
}
