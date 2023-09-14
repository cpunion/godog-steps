package file

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cpunion/godog-steps/cmd"
	"github.com/cucumber/godog"
	jd "github.com/nsf/jsondiff"
)

var assetsFS fs.FS

func InitAssetsFS(fs fs.FS) {
	assetsFS = fs
}

func iHaveAFileWithContent(ctx context.Context, file string, content string) (context.Context, error) {
	path := file
	if !strings.HasPrefix(file, "/") {
		path = filepath.Join(cmd.GetCwd(ctx), file)
	}
	os.MkdirAll(filepath.Dir(path), 0755)
	err := os.WriteFile(path, []byte(content), 0644)
	return ctx, err
}

func iHaveAFileFrom(ctx context.Context, file string, from string) (context.Context, error) {
	ffrom, err := assetsFS.Open(from)
	if err != nil {
		return ctx, err
	}
	defer ffrom.Close()
	path := file
	if !strings.HasPrefix(file, "/") {
		path = filepath.Join(cmd.GetCwd(ctx), file)
	}
	os.MkdirAll(filepath.Dir(path), 0755)
	fto, err := os.Create(path)
	if err != nil {
		return ctx, err
	}
	defer fto.Close()
	_, err = io.Copy(fto, ffrom)
	return ctx, err
}

func unindent(s string) (string, error) {
	s = strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return "", nil
	}
	// get indent of first line
	idtStr := ""
	for _, c := range lines[0] {
		if c != ' ' && c != '\t' {
			break
		}
		idtStr += string(c)
	}
	// trim indent of all lines
	for i, line := range lines {
		if strings.HasPrefix(line, idtStr) {
			lines[i] = line[len(idtStr):]
		} else if line != "" {
			return "", fmt.Errorf("line %d of %s has invalid indent", i+1, s)
		}
	}
	return strings.Join(lines, "\n"), nil
}

func iShouldHaveAFileWithContent(ctx context.Context, file string, content string) (context.Context, error) {
	content, err := unindent(content)
	if err != nil {
		return ctx, err
	}
	path := file
	if !strings.HasPrefix(file, "/") {
		path = filepath.Join(cmd.GetCwd(ctx), file)
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return ctx, err
	}
	fileContent := string(bytes)
	if fileContent != content {
		fmt.Printf("content: [%d]%v, got: [%d]%v", len(content), content, len(fileContent), fileContent)
		return ctx, fmt.Errorf("file content of %s mismatch, expected: %s, got: %s", file, content, string(bytes))
	}
	return ctx, nil
}

func iShouldHaveAJSONFileLikes(ctx context.Context, file, content string) (context.Context, error) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return ctx, err
	}
	diff, diffData := jd.Compare([]byte(content), fileContent, &jd.Options{})
	if diff != jd.FullMatch {
		return ctx, fmt.Errorf("file content of %s mismatch, expected: %s, got: %s, difference: %s", file, content, string(fileContent), string(diffData))
	}
	return ctx, nil
}

func iHaveASymbolLinkTo(ctx context.Context, link, path string) (context.Context, error) {
	if err := os.MkdirAll(filepath.Dir(link), 0755); err != nil {
		return ctx, err
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(cmd.GetCwd(ctx), path)
	}
	return ctx, os.Symlink(path, link)
}

func Init(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		serverState := &ctxServerState{mapping: make(map[string]string)}
		ctx = context.WithValue(ctx, ctxServerKey{}, serverState)
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		server := getServer(ctx)
		if server != nil {
			server.Shutdown(ctx)
		}
		return ctx, err
	})
	ctx.Step(`^I have a file "([^"]*)" with content:$`, iHaveAFileWithContent)
	ctx.Step(`^I have a file "([^"]*)" from "([^"]*)"$`, iHaveAFileFrom)
	ctx.Step(`^I have a symbol link "([^"]*)" to "([^"]*)"$`, iHaveASymbolLinkTo)
	ctx.Step(`^I have a url "([^"]*)" from file "([^"]*)"$`, iHaveAUrlFromFile)
	ctx.Step(`^I should have a file "([^"]*)" with content:$`, iShouldHaveAFileWithContent)
	ctx.Step(`^I should have a JSON file "([^"]*)" likes:$`, iShouldHaveAJSONFileLikes)
}
