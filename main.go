package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/sheepla/qiitaz/client"
	"github.com/sheepla/qiitaz/ui"
	"github.com/toqueteos/webbrowser"
)

// nolint:gochecknoglobals
var (
	appName     = "qiitaz"
	appVersion  = "unknown"
	appRevision = "unknown"
	appUsage    = "[OPTIONS] QUERY..."
)

type exitCode int

// nolint:maligned
type options struct {
	Version bool   `short:"V" long:"version" description:"Show version"`
	Sort    string `short:"s" long:"sort" description:"Sort key to search e.g. \"created\", \"like\", \"stock\", \"rel\",  (default: \"rel\")" `
	Open    bool   `short:"o" long:"open" description:"Open URL in your web browser"`
	Preview bool   `short:"p" long:"preview" description:"Preview page on your terminal"`
	PageNo  int    `short:"n" long:"pageno" description:"Max page number of search page" default:"1"`
	JSON    bool   `short:"j" long:"json" description:"Output result in JSON format"`
}

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrRequest
	exitCodeErrFuzzyFinder
	exitCodeErrWebbrowser
	exitCodeErrJSON
	exitCodeErrPreview
)

func main() {
	exitCode, err := Main(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(int(exitCode))
}

// nolint:funlen,golint,revive,cyclop
func Main(cliArgs []string) (exitCode, error) {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = appName
	parser.Usage = appUsage

	args, err := parser.ParseArgs(cliArgs)
	// The content of the go-flags error is already output, so ignore it.
	if err != nil {
		if flags.WroteHelp(err) {
			return exitCodeOK, nil
		}
		return exitCodeErrArgs, nil
	}

	if opts.Version {
		// nolint:forbidigo
		fmt.Printf("%s: v%s-%s\n", appName, appVersion, appRevision)

		return exitCodeOK, nil
	}

	if len(args) == 0 {
		// nolint:goerr113
		return exitCodeErrArgs, errors.New("must require argument (s)")
	}

	if opts.PageNo <= 0 {
		fmt.Fprintln(os.Stderr)
		// nolint:goerr113
		return exitCodeErrArgs, errors.New("the page number must be a positive value")
	}

	var urls []string

	for i := 1; i <= opts.PageNo; i++ {
		u, err := client.NewSearchURL(strings.Join(args, " "), client.SortBy(opts.Sort), i)
		if err != nil {
			return exitCodeErrArgs, fmt.Errorf("failed to create search URL %s: %w", u, err)
		}

		urls = append(urls, u)
	}

	var results []client.Result

	for _, u := range urls {
		r, err := client.Search(u)
		if err != nil {
			return exitCodeErrRequest, fmt.Errorf("failed to search articles: %w", err)
		}

		results = append(results, r...)
	}

	if len(results) == 0 {
		// nolint:goerr113
		return exitCodeOK, errors.New("no results found")
	}

	if opts.JSON {
		bytes, err := json.Marshal(&results)
		if err != nil {
			return exitCodeErrJSON, fmt.Errorf("failed to marshalling JSON: %w", err)
		}

		stdout := bufio.NewWriter(os.Stdout)
		fmt.Fprintln(stdout, string(bytes))
		stdout.Flush()
		return exitCodeOK, nil
	}

	if opts.Preview {
		if err := startPreviewMode(results); err != nil {
			return exitCodeErrPreview, fmt.Errorf("an error occurred on preview mode: %w", err)
		}
		return exitCodeOK, nil
	}

	choices, err := ui.FindMulti(results)
	if err != nil {
		return exitCodeErrFuzzyFinder, fmt.Errorf("an error occurred on fuzzyfinder: %w", err)
	}

	if len(choices) == 0 {
		return exitCodeOK, nil
	}

	if opts.Open {
		for _, idx := range choices {
			url := client.NewPageURL(results[idx].Link)
			if err := webbrowser.Open(url); err != nil {
				return exitCodeErrWebbrowser, fmt.Errorf("failed to open the URL %s: %w", url, err)
			}
		}
	}

	for _, idx := range choices {
		// nolint:forbidigo
		fmt.Println(client.NewPageURL(results[idx].Link))
	}

	return exitCodeOK, nil
}

func startPreviewMode(result []client.Result) error {
	for {
		idx, err := ui.Find(result)
		if err != nil {
			if errors.Is(fuzzyfinder.ErrAbort, err) {
				// normal termination
				return nil
			}
			return fmt.Errorf("an error occurred on fuzzyfinder: %w", err)
		}

		title := result[idx].Title
		path := result[idx].Link

		pager, err := ui.NewPagerProgram(path, title)
		if err != nil {
			return fmt.Errorf("failed to init pager program: %w", err)
		}

		if err := pager.Start(); err != nil {
			return fmt.Errorf("an error occurred on pager: %w", err)
		}
	}
}
