package main

import (
	"flag"
	"log"
	"os"

	"github.com/otiai10/gh-dependents/ghdeps"
)

var (
	verbose    bool
	tpl        string
	sortByStar bool
	page       int
	json       bool
)

func main() {
	flag.BoolVar(&verbose, "v", false, "Show verbose log")
	flag.BoolVar(&json, "json", false, "Output in JSON format (Alias of -t=json)")
	flag.StringVar(&tpl, "t", "", "Output template ('' = default, 'json')")
	flag.BoolVar(&sortByStar, "s", false, "Output with sorting by num of stars")
	flag.IntVar(&page, "p", 0, "Pages to crawl (0 == all)")
	flag.Parse()
	identity := flag.Arg(0)
	c := &ghdeps.Crawler{
		ServiceURL: ghdeps.GitHubBaseURL,
		Source:     ghdeps.CreateRepository(identity),
		Verbose:    verbose,
	}
	if err := c.Crawl(page); err != nil {
		log.Fatalln(err)
	}
	opt := &ghdeps.PrintOption{SortByStar: sortByStar}
	switch {
	case tpl == "json", json:
		opt.Template = ghdeps.JSONTemplate
	}
	if err := c.Print(os.Stdout, opt); err != nil {
		log.Fatalln(err)
	}
}
