package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	Header  string   `json:"header"`
	Title   string   `json:"title"`
	Link    string   `json:"link"`
	Snippet string   `json:"snippet"`
	Tags    []string `json:"tags"`
}

type SortBy string

func (s SortBy) validate() bool {
	switch s {
	case "like":
		return true
	case "stock":
		return true
	case "rel":
		return true
	case "created":
		return true
	default:
		return false
	}
}

func NewSearchURL(query string, sortby SortBy, pageno int) (string, error) {
	if sortby == "" {
		sortby = "rel"
	}

	// nolint:goerr113
	if !sortby.validate() {
		return "", fmt.Errorf("invalid sort key: %s", sortby)
	}

	// nolint:exhaustivestruct,exhaustruct,varnamelen
	u := &url.URL{
		Scheme: "https",
		Host:   "qiita.com",
		Path:   "search",
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("sort", string(sortby))
	q.Set("page", strconv.Itoa(pageno))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// nolint:gosec,noctx
func Search(url string) ([]Result, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the page %s: %w", url, err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var (
		results []Result
		result  Result
	)

	doc.Find("div.searchResult_main").Each(func(i int, div *goquery.Selection) {
		result.Header = div.Find("div.searchResult_header").Text()
		t := div.Find("h1.searchResult_itemTitle a")
		result.Title = t.Text()
		result.Link = t.AttrOr("href", "")
		result.Snippet = div.Find("div.searchResult_snippet").Text()
		result.Tags = nil
		div.Find("li.tagList_item a").Each(func(i int, a *goquery.Selection) {
			result.Tags = append(result.Tags, a.Text())
		})
		results = append(results, result)
	})

	return results, nil
}

// nolint:exhaustivestruct,exhaustruct,varnamelen
func NewPageURL(path string) string {
	u := &url.URL{
		Scheme: "https",
		Host:   "qiita.com",
		Path:   path,
	}

	return u.String()
}

func newPageMarkdownURL(path string) string {
	return NewPageURL(path) + ".md"
}

func FetchArticle(path string) (io.ReadCloser, error) {
	url := newPageMarkdownURL(path)
	// nolint:gosec,noctx
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the page (%s): %w", url, err)
	}

	if res.StatusCode != http.StatusOK {
		// nolint:goerr113
		return nil, fmt.Errorf("HTTP status error: %d %s", res.StatusCode, res.Status)
	}

	return res.Body, nil
}
