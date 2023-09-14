package cmd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cpunion/godog-steps/helper"
	"github.com/cucumber/godog"
)

type cmdStateKey struct{}

type cmdState struct {
	cmd    *exec.Cmd
	cwd    string
	bg     bool
	print  bool
	stdOut bytes.Buffer
	// stdErr bytes.Buffer
	err error
}

func getCmdState(ctx context.Context) *cmdState {
	return ctx.Value(cmdStateKey{}).(*cmdState)
}

func getCmd(ctx context.Context) *exec.Cmd {
	cmdState := getCmdState(ctx)
	if cmdState == nil {
		return nil
	}
	return cmdState.cmd
}

func GetCwd(ctx context.Context) string {
	cmdState := getCmdState(ctx)
	if cmdState == nil {
		return ""
	}
	return cmdState.cwd
}

func GetErr(ctx context.Context) error {
	cmdState := getCmdState(ctx)
	return cmdState.err
}

func GetOut(ctx context.Context) string {
	cmdState := getCmdState(ctx)
	if !cmdState.bg {
		out, err := cmdState.cmd.CombinedOutput()
		if err != nil {
			return ""
		}
		return string(out)
	}
	return cmdState.stdOut.String()
}

func showConsoleOutputAfterExit(ctx context.Context) (context.Context, error) {
	cmdState := getCmdState(ctx)
	cmdState.print = true
	return ctx, nil
}

func iRun(ctx context.Context, command string) (context.Context, error) {
	args := strings.Fields(command)
	if len(args) == 0 {
		return ctx, errors.New("invalid command: " + command)
	}
	cmdState := getCmdState(ctx)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = GetCwd(ctx)
	cmd.Stdout = bufio.NewWriter(&cmdState.stdOut)
	cmd.Stderr = bufio.NewWriter(&cmdState.stdOut)
	cmdState.cmd = cmd
	cmdState.bg = false
	cmdState.err = cmd.Run()
	return ctx, nil
}

func iRunInBackground(ctx context.Context, command string) (context.Context, error) {
	args := strings.Fields(command)
	if len(args) == 0 {
		return ctx, errors.New("invalid command: " + command)
	}
	cmdState := getCmdState(ctx)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = GetCwd(ctx)
	cmd.Stdout = &cmdState.stdOut
	cmd.Stderr = &cmdState.stdOut
	cmdState.cmd = cmd
	cmdState.bg = true
	cmdState.err = cmd.Start()
	time.Sleep(10 * time.Millisecond)
	return ctx, cmdState.err
}

func itShouldRunSuccessfuly(ctx context.Context) (context.Context, error) {
	err := GetErr(ctx)
	if err != nil {
		out := GetOut(ctx)
		fmt.Printf("err: %+v\n", err)
		fmt.Printf("out: %+v\n", out)
	}
	return ctx, err
}

func iShouldHaveAFileAt(ctx context.Context, path string) (context.Context, error) {
	st, err := os.Stat(filepath.Join(GetCwd(ctx), path))
	if err != nil {
		return ctx, err
	}
	if st.IsDir() {
		return ctx, errors.New("file is a directory")
	}
	return ctx, nil
}

func iShouldSeeOutput(ctx context.Context, content string) (context.Context, error) {
	var out string
	toCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	err := helper.WaitFor(toCtx, func() bool {
		out = GetOut(ctx)
		return strings.Contains(out, content)
	})
	if err != nil {
		return ctx, fmt.Errorf("expect including %q, got %q, err: %v", content, out, err)
	}
	return ctx, nil
}

func Init(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// create a random temporary directory and change to it
		dir, err := os.MkdirTemp("", "gltf")
		if err != nil {
			return ctx, err
		}
		err = os.Chdir(dir)
		if err != nil {
			return ctx, err
		}
		cmdState := &cmdState{cwd: dir}
		ctx = context.WithValue(ctx, cmdStateKey{}, cmdState)
		return ctx, err
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		cmdState := getCmdState(ctx)
		if cmdState != nil {
			os.RemoveAll(cmdState.cwd)
		}
		return ctx, err
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		cmdState := getCmdState(ctx)
		cmd := cmdState.cmd
		if cmd != nil && cmdState.bg {
			err := cmd.Process.Signal(os.Interrupt)
			if err != nil {
				fmt.Printf("failed to send interrupt signal to process: %v", err)
				return ctx, err
			}
			err = cmd.Wait()
			if cmdState.print {
				fmt.Printf("process exited with code %d", cmd.ProcessState.ExitCode())
				fmt.Printf("output:\n%s", cmdState.stdOut.String())
			}
			if err != nil {
				fmt.Printf("output:\n%s", cmdState.stdOut.String())
				return ctx, err
			}
		}
		if err != nil {
			fmt.Printf("output:\n%s", cmdState.stdOut.String())
			return ctx, err
		}
		return ctx, nil
	})

	ctx.Step(`^I run "([^"]*)"$`, iRun)
	ctx.Step(`^I run "([^"]*)" in background$`, iRunInBackground)
	ctx.Step(`^It should run successfully$`, itShouldRunSuccessfuly)
	ctx.Step(`^I should see output "([^"]*)"$`, iShouldSeeOutput)
	ctx.Step(`^I should have a file "([^"]*)"$`, iShouldHaveAFileAt)
	ctx.Step(`^Show console output after exit$`, showConsoleOutputAfterExit)
}
