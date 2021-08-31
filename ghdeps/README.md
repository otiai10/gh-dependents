# ghdeps

Web crawler package to collect dependents of a repository.

# Usage

```go
import (
    "github.com/otiai10/gh-dependents/ghdeps"
)

func main() {
    c := ghdeps.NewCrawler("otiai10/lookpath")
    if err := c.All(); err != nil {
        panic(err)
    }
    c.Print(os.Stdout, ghdeps.PrintOption{
        Template: ghdeps.JSONTemplate,
    })
}

```
