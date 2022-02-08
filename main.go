package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/sheepla/qiitaz/client"
)

func main() {
	url := client.NewURL(os.Args[1], "like")
	result, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	choices, err := find(result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("You selected %d\n", choices)
}

func find(result []client.Result) ([]int, error) {
	return fuzzyfinder.FindMulti(
		result,
		func(i int) string {
			return result[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			return fmt.Sprintf("%s\n\n%s\n\n%s", result[i].Header, result[i].Title, result[i].Snippet)
		}),
	)
}
