package client

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	Header  string `qoquery:"div.searchResult_header,text" json:"header"`
	Title   string `qoquery:"h1.searchResult_itemTitle a,text" json:"title"`
	Link    string `qoquery:"h1.searchResult_itemTitle a,[href]" json:"link"`
	Snippet string `qoquery:"div.searchResult_snippet,text" json:"snippet"`
	// Tags    []string `qoquery:"li.tagList_item a,text" json:"tags"`
	// Likes   string   `qoquery:"ul.searchResult_statusList li,text" json:"likes"`
}

type SortBy string

const (
	Created SortBy = "created"
	Like    SortBy = "like"
	Stock   SortBy = "stock"
	Rel     SortBy = "rel"
)

func NewURL(query string, sortby SortBy) string {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "qiita.com"
	u.Path = "search"
	q := u.Query()
	q.Set("q", query)
	q.Set("sort", string(sortby))
	u.RawQuery = q.Encode()
	return u.String()
}

func Get(url string) ([]Result, error) {
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
		results = append(results, r)
	})

	return results, err
}