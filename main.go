package main

import (
	"flag"
	"log"
	"os"

	"github.com/otiai10/gh-dependents/ghdeps"
)

var (
	verbose bool
	tpl     string
	sorter  string
	page    int
	after   string
	json    bool
	pretty  bool
)

func main() {
	flag.BoolVar(&verbose, "v", false, "Show verbose log")
	flag.BoolVar(&pretty, "pretty", false, "Output in pretty format (Alias of -t=pretty)")
	flag.BoolVar(&json, "json", false, "Output in JSON format (Alias of -t=json)")
	flag.StringVar(&tpl, "t", "", "Output template ('json' = default, 'pretty')")
	flag.StringVar(&sorter, "s", "", "Sort ('' = default, 'star', 'fork')")
	flag.IntVar(&page, "p", 0, "Pages to crawl (0 == all)")
	flag.StringVar(&after, "a", "", "Hash of offset to be set in `dependents_after` query param")
	flag.Parse()
	identity := flag.Arg(0)
	c := &ghdeps.Crawler{
		ServiceURL: ghdeps.GitHubBaseURL,
		Source:     ghdeps.CreateRepository(identity),
		Verbose:    verbose,
		After:      after,
		PageCount:  page,
	}
	if err := c.Crawl(); err != nil {
		log.Fatalln(err)
	}
	opt := &ghdeps.PrintOption{}
	switch {
	case tpl == "pretty", pretty:
		opt.Template = ghdeps.PrettyTemplate
	case tpl == "json", json:
		opt.Template = ghdeps.JSONTemplate
	default:
		opt.Template = ghdeps.JSONTemplate
	}

	switch sorter {
	case "star", "stars":
		opt.Sort = ghdeps.SortByStar
	case "fork", "forks":
		opt.Sort = ghdeps.SortByFork
	}
	if err := c.Print(os.Stdout, opt); err != nil {
		log.Fatalln(err)
	}
}
