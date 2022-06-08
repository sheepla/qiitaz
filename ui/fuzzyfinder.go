package ui

import (
	"fmt"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
	"github.com/sheepla/qiitaz/client"
)

// nolint:wrapcheck,gomnd
func FindMulti(result []client.Result) ([]int, error) {
	return fuzzyfinder.FindMulti(
		result,
		func(i int) string {
			return result[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(idx, width, height int) string {
			if idx == -1 {
				return ""
			}

			wrapedWidth := width/2 - 5

			return runewidth.Wrap(renderPreviewWindow(&result[idx]), wrapedWidth)
		}),
	)
}

// nolint:wrapcheck,gomnd
func Find(result []client.Result) (int, error) {
	return fuzzyfinder.Find(
		result,
		func(i int) string {
			return result[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(idx, width, height int) string {
			if idx == -1 {
				return ""
			}

			wrapedWidth := width/2 - 5

			return runewidth.Wrap(renderPreviewWindow(&result[idx]), wrapedWidth)
		}),
	)
}

func renderPreviewWindow(result *client.Result) string {
	return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
		result.Header,
		result.Title,
		result.Snippet,
		strings.Join(result.Tags, " "),
	)
}
