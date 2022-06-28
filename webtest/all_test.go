package webtest

import (
	"testing"

	"github.com/otiai10/gh-dependents/ghdeps"
	. "github.com/otiai10/mint"
)

func TestCrawler_All(t *testing.T) {
	c := ghdeps.NewCrawler("cli/cli")
	err := c.All()
	Expect(t, err).ToBe(nil)
	Expect(t, len(c.Pages) > 0).ToBe(true)
	Expect(t, len(c.Dependents) >= 20).ToBe(true)
}

func TestCrawler_Crawl(t *testing.T) {
	c := ghdeps.NewCrawler("otiai10/lookpath")
	c.PageCount = 1
	err := c.Crawl()
	Expect(t, err).ToBe(nil)
	Expect(t, len(c.Pages)).ToBe(1)
	Expect(t, len(c.Dependents)).ToBe(29)
}
