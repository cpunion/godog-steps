package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cpunion/godog-steps/cmd"
	"github.com/cucumber/godog"
	"github.com/kettek/apng"
)

func theFileShouldBeAnAnimatedPng(ctx context.Context, file string) (context.Context, error) {
	cwd := cmd.GetCwd(ctx)
	r, err := os.Open(filepath.Join(cwd, file))
	if err != nil {
		return ctx, err
	}
	defer r.Close()
	a, err := apng.DecodeAll(r)
	if err != nil {
		return ctx, err
	}
	if len(a.Frames) <= 1 {
		return ctx, fmt.Errorf("file %q is not an animated png, frames: %d", file, len(a.Frames))
	}
	return ctx, nil
}

func Init(ctx *godog.ScenarioContext) {
	ctx.Step(`^the file "([^"]*)" should be an animated png$`, theFileShouldBeAnAnimatedPng)
}
