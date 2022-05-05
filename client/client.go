package client

import (
	"fmt"
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
	if !sortby.validate() {
		return "", fmt.Errorf("invalid sort key: %s", sortby)
	}
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

func Search(url string) ([]Result, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the page %s: %s", url, err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %s", err)
	}

	var (
		results []Result
		r       Result
	)
	doc.Find("div.searchResult_main").Each(func(i int, div *goquery.Selection) {
		r.Header = div.Find("div.searchResult_header").Text()
		t := div.Find("h1.searchResult_itemTitle a")
		r.Title = t.Text()
		r.Link = t.AttrOr("href", "")
		r.Snippet = div.Find("div.searchResult_snippet").Text()
		r.Tags = nil
		div.Find("li.tagList_item a").Each(func(i int, a *goquery.Selection) {
			r.Tags = append(r.Tags, a.Text())
		})
		results = append(results, r)
	})

	return results, nil
}

func NewPageURL(path string) string {
	u := &url.URL{
		Scheme: "https",
		Host:   "qiita.com",
		Path:   path,
	}
	return u.String()
}

func NewPageMarkdownURL(path string) string {
	return NewPageURL(path) + ".md"
}
