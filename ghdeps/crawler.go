package ghdeps

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	GitHubBaseURL             = "https://github.com"
	defaultSleepIntervalPages = 12
)

type (
	Crawler struct {
		ServiceURL string
		Source     Repository
		Dependents Dependents
		Pages      []Page

		// Configs
		After              string
		PageNum            int
		Verbose            bool
		SleepIntervalPages int
	}

	Dependents []Repository

	Repository struct {
		User  string
		Repo  string
		Stars int
	}

	Page struct {
		URL  string
		Next string
		// YAGNI: Previous string
	}
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

func NewCrawler(identity string) *Crawler {
	return &Crawler{
		ServiceURL: GitHubBaseURL,
		Source:     CreateRepository(identity),
	}
}

func (c *Crawler) All() error {
	return c.Crawl(0)
}

func (c *Crawler) Crawl(page int) (err error) {
	c.PageNum = page
	if c.SleepIntervalPages == 0 {
		c.SleepIntervalPages = defaultSleepIntervalPages
	}
	link := c.Source.URL(c.ServiceURL) + "/network/dependents"
	if c.After != "" {
		link += "?dependents_after=" + c.After
	}
	for link != "" {
		c.sleepIfNeeded()
		if link, err = c.Page(link); err != nil {
			return err
		}
		if page != 0 && len(c.Pages) >= page {
			return nil
		}
	}
	return nil
}

func (c *Crawler) sleepIfNeeded() {
	if len(c.Pages) != 0 && len(c.Pages)%c.SleepIntervalPages == 0 {
		rand.Seed(time.Now().Unix())
		secs := rand.Intn(12)
		if len(c.Pages)%(c.SleepIntervalPages*5) == 0 {
			secs = secs * 5
		}
		if c.Verbose {
			fmt.Fprintf(os.Stderr, "Sleeping %d seconds to avoid HTTP 429.\n", secs)
		}
		time.Sleep(time.Duration(secs) * time.Second)
	}
}

func (c *Crawler) Page(link string) (next string, err error) {
	if c.Verbose {
		fmt.Fprintf(os.Stderr, "[Page % 2d] %s", len(c.Pages)+1, link)
	}
	res, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf(res.Status)
	}
	next, err = c.page(res.Body)
	if c.Verbose {
		fmt.Fprintf(os.Stderr, "\t= %d\n", len(c.Dependents))
	}
	u, _ := url.Parse(next)
	c.Pages = append(c.Pages, Page{URL: link, Next: u.Query().Get("dependents_after")})
	return next, err
}

func (c *Crawler) page(r io.Reader) (next string, err error) {
	node, err := html.Parse(r)
	if err != nil {
		return "", err
	}
	if c.Dependents == nil {
		c.Dependents = Dependents{}
	}
	next, err = c.Walk(node)
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

func (r Repository) URL(base string) string {
	return fmt.Sprintf("%s/%s/%s", base, r.User, r.Repo)
}

func CreateRepository(identifier string) Repository {
	id := strings.Split(strings.Trim(identifier, "/"), "/")
	if len(id) < 2 {
		log.Fatalf("Failed to parse repository identity: %s\n", identifier)
	}
	return Repository{User: id[0], Repo: id[1]}
}

func CreateRepositoryFromRowNode(node *html.Node) (repo Repository, err error) {
	a := node.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.NextSibling.NextSibling.NextSibling
	repo = CreateRepository(getAttribute(a, "href"))
	stars := node.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling
	numstars, err := strconv.Atoi(noiseOfStars.ReplaceAllString(stars.Data, ""))
	if err != nil {
		numstars = 0 // TODO: Does GitHub use like "1M" for "1,000,000"?
	}
	repo.Stars = numstars
	return repo, nil
}
