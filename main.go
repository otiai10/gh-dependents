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
)

func main() {
	flag.BoolVar(&verbose, "v", false, "Show verbose log")
	flag.StringVar(&tpl, "tpl", "", "Output template ('' = default, 'json')")
	flag.Parse()
	identity := flag.Arg(0)
	c := &ghdeps.Crawler{
		ServiceURL: ghdeps.GitHubBaseURL,
		Source:     ghdeps.CreateRepository(identity),
		Verbose:    verbose,
	}
	if err := c.All(); err != nil {
		log.Fatalln(err)
	}
	opt := new(ghdeps.PrintOption)
	switch tpl {
	case "json":
		opt.Template = ghdeps.JSONTemplate
	}
	if err := c.Print(os.Stdout, opt); err != nil {
		log.Fatalln(err)
	}
}
