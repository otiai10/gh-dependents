package ghdeps

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

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

func (c *Crawler) All() (err error) {
	link := c.Source.URL(c.ServiceURL) + "/network/dependents"
	for link != "" {
		if link, err = c.Page(link); err != nil {
			return err
		}
	}
	sort.Sort(c.Dependents)
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

func (c *Crawler) Walk(node *html.Node) (string, error) {
	if node.Type == html.ElementNode {
		if node.Data == "div" {
			if getAttribute(node, "id") == "dependents" {
				// assume 4th child is the box
				box := node.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling
				// assume 1st child of the box is header
				for row := box.FirstChild.NextSibling.NextSibling.NextSibling; row != nil; row = row.NextSibling {
					if row.Type != html.ElementNode || row.Data != "div" {
						continue
					}
					repo, err := CreateRepositoryFromRowNode(row)
					if err != nil {
						return "", err
					}
					c.Dependents = append(c.Dependents, repo)
				}
				page := box.NextSibling.NextSibling.FirstChild.NextSibling
				for btn := page.FirstChild; btn != nil; btn = btn.NextSibling {
					if btn.Type == html.ElementNode && btn.Data == "a" {
						href := getAttribute(btn, "href")
						if strings.Contains(href, "dependents_after") {
							return href, nil
						}
					}
				}
				return "", nil
			}
		}
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
	return opt.Template.Execute(out, c)
}

func getAttribute(node *html.Node, name string) string {
	if node == nil || node.Attr == nil {
		return ""
	}
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}
