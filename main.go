package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/sheepla/qiitaz/client"
	"github.com/toqueteos/webbrowser"
)

const (
	appName    = "qiitaz"
	appVersion = "0.0.3"
	appUsage   = "[OPTIONS] QUERY..."
)

const (
	baseURL = "https://qiita.com"
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
		log.Println(err)
		return exitCodeErrArgs
	}

	result, err := client.Search(url)
	if err != nil {
		log.Println(err)
		return exitCodeErrRequest
	}

	choices, err := find(result)
	if err != nil {
		log.Println(err)
		return exitCodeErrFuzzyFinder
	}

	if opts.Open {
		for _, idx := range choices {
			url := path.Join(baseURL, result[idx].Link)
			if err := webbrowser.Open(url); err != nil {
				log.Println(err)
				return exitCodeErrWebbrowser
			}
		}
	}

	if opts.Preview {
		for _, idx := range choices {
			url := client.NewPageURL((result[idx].Link + ".md"))
			view, err := client.Preview(url, "dark")
			if err != nil {
				log.Println(err)
				return exitCodeErrPreview
			}
			fmt.Fprintf(os.Stdout, *view)
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
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
				result[i].Header,
				result[i].Title,
				result[i].Snippet,
				strings.Join(result[i].Tags, " "),
			)
		}),
	)
}
