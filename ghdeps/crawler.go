package ghdeps

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"

	"golang.org/x/net/html"
)

const GitHubBaseURL = "https://github.com"

type (
	Crawler struct {
		ServiceURL string
		Source     Repository
		Dependents Dependents
		Pages      []string

		// Configs
		Verbose bool
	}
	Dependents []Repository
)

var (
	noiseOfStars = regexp.MustCompile("[ \\n,]")
)

func (deps Dependents) Len() int {
	return len(deps)
}

func (deps Dependents) Less(i, j int) bool {
	return deps[i].Stars > deps[j].Stars
}

func (deps Dependents) Swap(i, j int) {
	deps[i], deps[j] = deps[j], deps[i]
}

func (c *Crawler) Crawl(page int) (err error) {
	link := c.Source.URL(c.ServiceURL) + "/network/dependents"
	for link != "" {
		if link, err = c.Page(link); err != nil {
			return err
		}
		if page != 0 && len(c.Pages) >= page {
			return nil
		}
	}
	return nil
}

func (c *Crawler) Page(link string) (string, error) {
	c.Pages = append(c.Pages, link)
	if c.Verbose {
		fmt.Fprintf(os.Stderr, "[Page % 2d] %s", len(c.Pages), link)
	}
	res, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf(res.Status)
	}
	next, err := c.page(res.Body)
	if c.Verbose {
		fmt.Fprintf(os.Stderr, "\t= %d\n", len(c.Dependents))
	}
	return next, err
}

func (c *Crawler) page(r io.Reader) (string, error) {
	node, err := html.Parse(r)
	if err != nil {
		return "", err
	}
	if c.Dependents == nil {
		c.Dependents = Dependents{}
	}
	next, err := c.Walk(node)
	if err != nil {
		return "", err
	}
	if next != "" {
		return next, nil
	}
	return "", nil
}

// Walk walkthrough all DOM element on the page recursively.
// See `query.go` for how we find the target nodes from DOM structure.
func (c *Crawler) Walk(node *html.Node) (string, error) {
	if box := queryBox(node); box != nil {
		for row := queryNextRow(box.FirstChild); row != nil; row = queryNextRow(row.NextSibling) {
			repo, err := CreateRepositoryFromRowNode(row)
			if err != nil {
				return "", err
			}
			c.Dependents = append(c.Dependents, repo)
		}
		if btn := queryNextPageButton(box); btn != nil {
			return getAttribute(btn, "href"), nil
		}
		return "", nil
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		next, err := c.Walk(child)
		if err != nil {
			return "", err
		}
		if next != "" {
			return next, nil
		}
	}
	return "", nil
}

func (c *Crawler) Print(out io.Writer, opt *PrintOption) error {
	opt = opt.ensure()
	if opt.SortByStar {
		sort.Sort(c.Dependents)
	}
	return opt.Template.Execute(out, c)
}
