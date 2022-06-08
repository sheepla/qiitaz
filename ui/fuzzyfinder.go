package ui

import (
	"fmt"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
	"github.com/sheepla/qiitaz/client"
)

func FindMulti(result []client.Result) ([]int, error) {
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
			return runewidth.Wrap(renderPreviewWindow(&result[i]), wrapedWidth)
		}),
	)
}

func Find(result []client.Result) (int, error) {
	return fuzzyfinder.Find(
		result,
		func(i int) string {
			return result[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			wrapedWidth := w/2 - 5
			return runewidth.Wrap(renderPreviewWindow(&result[i]), wrapedWidth)
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
