package ui

import (
	"fmt"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-runewidth"
	"github.com/sheepla/qiitaz/client"
)

func Find(result []client.Result) ([]int, error) {
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

			return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
				runewidth.Wrap(result[i].Header, wrapedWidth),
				runewidth.Wrap(result[i].Title, wrapedWidth),
				runewidth.Wrap(result[i].Snippet, wrapedWidth),
				runewidth.Wrap(strings.Join(result[i].Tags, " "), wrapedWidth),
			)
		}),
	)
}
