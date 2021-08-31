package ghdeps

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/otiai10/mint"
)

func TestCrawler_Page(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo/baa/network/dependents", func(w http.ResponseWriter, req *http.Request) {
		f, _ := os.Open("./testdata/example.html")
		io.Copy(w, f)
	})
	server := httptest.NewServer(mux)
	source := CreateRepository("foo/baa")
	c := &Crawler{ServiceURL: server.URL, Source: source}
	next, err := c.Page(source.URL(c.ServiceURL) + "/network/dependents")
	Expect(t, err).ToBe(nil)
	Expect(t, next).ToBe("https://github.com/otiai10/copy/network/dependents?dependents_after=MTY3MjUwNTM5MzU")

	buf := bytes.NewBuffer(nil)
	err = c.Print(buf, nil)
	Expect(t, err).ToBe(nil)
	Expect(t, buf.String()).Match("TOTAL:\t30")

	buf = bytes.NewBuffer(nil)
	err = c.Print(buf, &PrintOption{})
	Expect(t, err).ToBe(nil)
	Expect(t, buf.String()).Match("TOTAL:\t30")
}

func TestCrawler_All(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo/baa/network/dependents", func(w http.ResponseWriter, req *http.Request) {
		f, _ := os.Open("./testdata/no-next.html")
		io.Copy(w, f)
	})
	server := httptest.NewServer(mux)
	source := CreateRepository("foo/baa")
	c := &Crawler{ServiceURL: server.URL, Source: source}
	err := c.Crawl(0)
	Expect(t, err).ToBe(nil)
	Expect(t, len(c.Pages)).ToBe(1)
	Expect(t, len(c.Dependents)).ToBe(4)
}

func TestJSONTemplate(t *testing.T) {
	c := &Crawler{
		ServiceURL: "http://localhost:8080",
		Source:     Repository{User: "otiai10", Repo: "gh-dependents"},
		Pages:      []string{},
		Dependents: []Repository{
			{User: "foo", Repo: "baa"},
		},
	}
	buf := bytes.NewBuffer(nil)
	err := c.Print(buf, &PrintOption{Template: JSONTemplate, SortByStar: true})
	Expect(t, err).ToBe(nil)
	out := map[string]interface{}{}
	err = json.NewDecoder(buf).Decode(&out)
	Expect(t, err).ToBe(nil)
}
