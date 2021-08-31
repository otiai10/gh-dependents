package ghdeps

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type (
	Repository struct {
		User  string
		Repo  string
		Stars int
	}
)

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
