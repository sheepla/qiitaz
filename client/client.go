package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	Header  string   `qoquery:"div.searchResult_header,text" json:"header"`
	Title   string   `qoquery:"h1.searchResult_itemTitle a,text" json:"title"`
	Link    string   `qoquery:"h1.searchResult_itemTitle a,[href]" json:"link"`
	Snippet string   `qoquery:"div.searchResult_snippet,text" json:"snippet"`
	Tags    []string `qoquery:"li.tagList_item a,text" json:"tags"`
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
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var results []Result
	var r Result

	doc.Find("div.searchResult_main").Each(func(i int, s *goquery.Selection) {
		r.Header = s.Find("div.searchResult_header").Text()
		t := s.Find("h1.searchResult_itemTitle a")
		r.Title = t.Text()
		r.Link = t.AttrOr("href", "")
		r.Snippet = s.Find("div.searchResult_snippet").Text()
		r.Tags = nil
		s.Find("li.tagList_item a").Each(func(i int, ss *goquery.Selection) {
			r.Tags = append(r.Tags, ss.Text())
		})
		results = append(results, r)
	})

	return results, err
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
