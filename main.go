package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
	"github.com/sheepla/qiitaz/client"
	"github.com/sheepla/qiitaz/ui"
	"github.com/toqueteos/webbrowser"
)

const (
	appName    = "qiitaz"
	appVersion = "0.0.5"
	appUsage   = "[OPTIONS] QUERY..."
)

type exitCode int

type options struct {
	Version bool   `short:"V" long:"version" description:"Show version"`
	Sort    string `short:"s" long:"sort" description:"Sort key to search e.g. \"created\", \"like\", \"stock\", \"rel\",  (default: \"rel\")" `
	Open    bool   `short:"o" long:"open" description:"Open URL in your web browser"`
	Preview bool   `short:"p" long:"preview" description:"Preview page on your terminal"`
	PageNo  int    `short:"n" long:"pageno" description:"Number of search page"`
}

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrRequest
	exitCodeErrFuzzyFinder
	exitCodeErrWebbrowser
	exitCodeErrPreview
)

func main() {
	os.Exit(int(Main(os.Args[1:])))
}

func Main(cliArgs []string) exitCode {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = appName
	parser.Usage = appUsage

	args, err := parser.ParseArgs(cliArgs)
	if err != nil {
		if flags.WroteHelp(err) {
			return exitCodeOK
		} else {
			fmt.Fprintf(os.Stderr, "Argument parsing failed: %s", err)
			return exitCodeErrArgs
		}
	}

	if opts.Version {
		fmt.Printf("%s: v%s\n", appName, appVersion)
		return exitCodeOK
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Must require argument (s)")
		return exitCodeErrArgs
	}

	url, err := client.NewSearchURL(strings.Join(args, " "), client.SortBy(opts.Sort), opts.PageNo)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErrArgs
	}

	result, err := client.Search(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErrRequest
	}

	choices, err := find(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErrFuzzyFinder
	}

	if opts.Open {
		for _, idx := range choices {
			url := client.NewPageURL(result[idx].Link)
			if err := webbrowser.Open(url); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return exitCodeErrWebbrowser
			}
		}
	}

	if opts.Preview {
		for _, idx := range choices {
			url := client.NewPageURL((result[idx].Link + ".md"))
			title := result[idx].Title
			if err := ui.Preview(url, title); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return exitCodeErrPreview
			}
		}
	}

	return exitCodeOK
}

func find(result []client.Result) ([]int, error) {
	return fuzzyfinder.FindMulti(
		result,
		func(i int) string {
			return result[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			wrapedWidth := w/2 - 5

			header := runewidth.Wrap(result[i].Header, wrapedWidth)
			title := runewidth.Wrap(result[i].Title, wrapedWidth)
			snippet := runewidth.Wrap(result[i].Snippet, wrapedWidth)
			tags := runewidth.Wrap(strings.Join(result[i].Tags, " "), wrapedWidth)

			return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
				header,
				title,
				snippet,
				tags,
			)
		}),
	)
}
