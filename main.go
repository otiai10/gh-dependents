package main

import (
	"flag"
	"log"
	"os"

	"github.com/otiai10/gh-dependents/ghdeps"
)

var (
	verbose bool = false
)

func main() {
	flag.BoolVar(&verbose, "v", false, "Show verbose log")
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
	if err := c.Print(os.Stdout, nil); err != nil {
		log.Fatalln(err)
	}
}
