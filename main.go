package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/sheepla/qiitaz/client"
	"github.com/toqueteos/webbrowser"
)

const (
	baseURL = "https://qiita.com"
)

func main() {
	url := client.NewURL(os.Args[1], "like")
	result, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}

	choices, err := find(result)
	if err != nil {
		log.Println(err)
	}

	for _, idx := range choices {
		url := path.Join(baseURL, result[idx].Link)
		if err := webbrowser.Open(url); err != nil {
			log.Println(err)
		} else {
			fmt.Println(url)
		}
	}
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
			return fmt.Sprintf("%s\n\n%s\n\n%s", result[i].Header, result[i].Title, result[i].Snippet)
		}),
	)
}
