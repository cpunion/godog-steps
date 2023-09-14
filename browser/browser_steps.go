package browser

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/cucumber/godog"
)

type browserContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type browserKey struct{}

func getBrowser(ctx context.Context) *browserContext {
	browser := ctx.Value(browserKey{})
	if browser == nil {
		return nil
	}
	return browser.(*browserContext)
}

func setBrowser(ctx context.Context, cancel context.CancelFunc) context.Context {
	return context.WithValue(ctx, browserKey{}, &browserContext{ctx: ctx, cancel: cancel})
}

func iOpenTheBrowserTo(ctx context.Context, url string) (context.Context, error) {
	// remove Headless
	options := append(chromedp.DefaultExecAllocatorOptions[0:2], chromedp.DefaultExecAllocatorOptions[3:]...)
	chromeOpts := append(
		options,
		chromedp.Headless,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("use-angle", true),
	)

	allocContext, cancel := chromedp.NewExecAllocator(ctx, chromeOpts...)
	// create context
	bctx, _ := chromedp.NewContext(
		allocContext,
		// chromedp.WithDebugf(log.Printf),
		chromedp.WithLogf(log.Printf),
		chromedp.WithErrorf(log.Printf),
	)
	ctx = setBrowser(bctx, cancel)
	err := chromedp.Run(bctx, chromedp.Navigate(url))
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func iShouldSeeElement(ctx context.Context, node string) (context.Context, error) {
	browser := getBrowser(ctx)
	toCtx, cancel := context.WithTimeout(browser.ctx, 1*time.Second)
	defer cancel()
	err := chromedp.Run(toCtx, chromedp.WaitVisible(node, chromedp.ByQuery))
	return ctx, err
}

func iShouldSeeText(ctx context.Context, text string) (context.Context, error) {
	browser := getBrowser(ctx)
	toCtx, cancel := context.WithTimeout(browser.ctx, 1*time.Second)
	defer cancel()
	query := fmt.Sprintf("/"+"/*[contains(text(), '%s')]", strings.ReplaceAll(text, "'", "\\'"))
	err := chromedp.Run(toCtx, chromedp.WaitVisible(query, chromedp.BySearch))
	return ctx, err
}

func iClickText(ctx context.Context, text string) (context.Context, error) {
	browser := getBrowser(ctx)
	toCtx, cancel := context.WithTimeout(browser.ctx, 1*time.Second)
	defer cancel()
	query := fmt.Sprintf("/"+"/*[contains(text(), '%s')]", strings.ReplaceAll(text, "'", "\\'"))
	err := chromedp.Run(toCtx, chromedp.Click(query, chromedp.BySearch))
	return ctx, err
}

func iClick(ctx context.Context, query string) (context.Context, error) {
	browser := getBrowser(ctx)
	toCtx, cancel := context.WithTimeout(browser.ctx, 1*time.Second)
	defer cancel()
	err := chromedp.Run(toCtx, chromedp.Click(query, chromedp.ByQuery))
	return ctx, err
}

func iShouldSeeIs(ctx context.Context, query string, text string) (context.Context, error) {
	browser := getBrowser(ctx)
	toCtx, cancel := context.WithTimeout(browser.ctx, 1*time.Second)
	defer cancel()
	var value string
	err := chromedp.Run(toCtx, chromedp.Value(query, &value, chromedp.ByQuery))
	if err != nil {
		return ctx, err
	}
	if value != text {
		return ctx, fmt.Errorf("expected %q, got %q", text, value)
	}
	return ctx, nil
}

func iReplaceFormValuesWith(ctx context.Context, from, to string, table *godog.Table) (context.Context, error) {
	if from == to {
		return ctx, fmt.Errorf("from '%s' and to '%s' cannot be the same", from, to)
	}
	browser := getBrowser(ctx)
	toCtx, cancel := context.WithTimeout(browser.ctx, 1*time.Second)
	defer cancel()
	if len(table.Rows) < 2 {
		return ctx, fmt.Errorf("expected at least 2 rows, including header and rows, got %d", len(table.Rows))
	}
	// find from and to index
	fromIndex, toIndex := -1, -1
	for i, cell := range table.Rows[0].Cells {
		if cell.Value == from {
			fromIndex = i
		}
		if cell.Value == to {
			toIndex = i
		}
	}
	if fromIndex == -1 {
		return ctx, fmt.Errorf("expected to find %q in header row", from)
	}
	if toIndex == -1 {
		return ctx, fmt.Errorf("expected to find %q in header row", to)
	}
	maxIdx := fromIndex
	if toIndex > maxIdx {
		maxIdx = toIndex
	}
	// process values in table
	for idx, row := range table.Rows {
		if idx == 0 {
			continue
		}
		if len(row.Cells) < maxIdx+1 {
			return ctx, fmt.Errorf("expected at least %d cells, got %d", maxIdx+1, len(row.Cells))
		}
		fromValue := row.Cells[fromIndex].Value
		toValue := row.Cells[toIndex].Value
		// wait for visible
		query := fmt.Sprintf("/"+"/*[@value='%s']", strings.ReplaceAll(fromValue, "'", "\\'"))
		err := chromedp.Run(toCtx, chromedp.WaitVisible(query, chromedp.BySearch))
		if err != nil {
			return ctx, fmt.Errorf("error wait visible %s: %w", query, err)
		}
		// clear first and then send keys to avoid without onchange event
		err = chromedp.Run(toCtx, chromedp.SetValue(query, "", chromedp.BySearch))
		if err != nil {
			return ctx, fmt.Errorf("error clear value %s: %w", query, err)
		}
		err = chromedp.Run(toCtx, chromedp.SendKeys(query, toValue))
		if err != nil {
			return ctx, fmt.Errorf("error replace value from %s to %s: %w", fromValue, toValue, err)
		}
	}
	return ctx, nil
}

func Init(ctx *godog.ScenarioContext) (*godog.ScenarioContext, error) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		browser := getBrowser(ctx)
		if browser != nil {
			browser.cancel()
		}
		return ctx, nil
	})
	ctx.Step(`^I open the browser to "([^"]*)"$`, iOpenTheBrowserTo)
	ctx.Step(`^I should see "([^"]*)" element$`, iShouldSeeElement)
	ctx.Step(`^I should see text "([^"]*)"$`, iShouldSeeText)
	ctx.Step(`^I click text "([^"]*)"$`, iClickText)
	ctx.Step(`^I click "([^"]*)"$`, iClick)
	ctx.Step(`^I should see "([^"]*)" is "([^"]*)"$`, iShouldSeeIs)
	ctx.Step(`^I replace form values "([^"]*)" with "([^"]*)":$`, iReplaceFormValuesWith)
	return ctx, nil
}
