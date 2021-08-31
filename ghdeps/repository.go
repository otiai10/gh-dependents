package ghdeps

import (
	"fmt"
	"log"
	"strings"
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
