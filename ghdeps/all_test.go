package ghdeps

import (
	"bytes"
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
}
