package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/sheepla/qiitaz/client"
	"github.com/sheepla/qiitaz/ui"
	"github.com/toqueteos/webbrowser"
)

var (
	appName     = "qiitaz"
	appVersion  = "unknown"
	appRevision = "unknown"
	appUsage    = "[OPTIONS] QUERY..."
)

type exitCode int

type options struct {
	Version bool   `short:"V" long:"version" description:"Show version"`
	Sort    string `short:"s" long:"sort" description:"Sort key to search e.g. \"created\", \"like\", \"stock\", \"rel\",  (default: \"rel\")" `
	Open    bool   `short:"o" long:"open" description:"Open URL in your web browser"`
	Preview bool   `short:"p" long:"preview" description:"Preview page on your terminal"`
	PageNo  int    `short:"n" long:"pageno" description:"Max page number of search page" default:"1"`
	Json    bool   `short:"j" long:"json" description:"Output result in JSON format"`
}

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrRequest
	exitCodeErrFuzzyFinder
	exitCodeErrWebbrowser
	exitCodeErrJson
	exitCodeErrPreview
)

func main() {
	exitCode, err := Main(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(int(exitCode))
}

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
		fmt.Printf("%s: v%s\n", appName, appVersion)
		return exitCodeOK, nil
	}

	if len(args) == 0 {
		return exitCodeErrArgs, errors.New("must require argument (s)")
	}

	if opts.PageNo <= 0 {
		fmt.Fprintln(os.Stderr)
		return exitCodeErrArgs, errors.New("the page number must be a positive value")
	}

	var urls []string
	for i := 1; i <= opts.PageNo; i++ {
		u, err := client.NewSearchURL(strings.Join(args, " "), client.SortBy(opts.Sort), i)
		if err != nil {
			return exitCodeErrArgs, fmt.Errorf("failed to create search URL %s: %s", u, err)
		}
		urls = append(urls, u)
	}

	var results []client.Result
	for _, u := range urls {
		r, err := client.Search(u)
		if err != nil {
			return exitCodeErrRequest, fmt.Errorf("failed to search articles: %s", err)
		}
		results = append(results, r...)
	}

	if len(results) == 0 {
		return exitCodeOK, errors.New("no results found")
	}

	if opts.Json {
		bytes, err := json.Marshal(&results)
		if err != nil {
			return exitCodeErrJson, fmt.Errorf("failed to marshalling JSON: %s", err)
		}
		stdout := bufio.NewWriter(os.Stdout)
		fmt.Fprintln(stdout, string(bytes))
		stdout.Flush()
		return exitCodeOK, nil
	}

	choices, err := ui.Find(results)
	if err != nil {
		return exitCodeErrFuzzyFinder, fmt.Errorf("an error occured on fuzzyfinder: %s", err)
	}

	if len(choices) == 0 {
		return exitCodeOK, nil
	}

	if opts.Open {
		for _, idx := range choices {
			url := client.NewPageURL(results[idx].Link)
			if err := webbrowser.Open(url); err != nil {
				return exitCodeErrWebbrowser, fmt.Errorf("failed to open the URL %s: %s", url, err)
			}
		}
	}

	if opts.Preview {
		for _, idx := range choices {
			url := client.NewPageMarkdownURL(results[idx].Link)
			title := results[idx].Title
			if err := ui.Preview(url, title); err != nil {
				return exitCodeErrPreview, fmt.Errorf("failed to preview the page (URL: %s, title: %s): %s", url, title, err)
			}
		}
	}

	for _, idx := range choices {
		fmt.Println(client.NewPageURL(results[idx].Link))
	}

	return exitCodeOK, nil
}
