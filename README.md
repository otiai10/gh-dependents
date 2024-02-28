# gh-dependents

[![Actions Status](https://github.com/otiai10/gh-dependents/workflows/Go/badge.svg)](https://github.com/otiai10/gh-dependents/actions)
[![Web Test](https://github.com/otiai10/gh-dependents/actions/workflows/webtest.yaml/badge.svg)](https://github.com/otiai10/gh-dependents/actions/workflows/webtest.yaml)
[![codecov](https://codecov.io/gh/otiai10/gh-dependents/branch/main/graph/badge.svg)](https://codecov.io/gh/otiai10/gh-dependents)

`gh` command extension to see dependents of your repository.

![screenshot](https://raw.githubusercontent.com/otiai10/gh-dependents/main/screenshot.png)

See The GitHub Blog: [GitHub CLI 2.0 includes extensions!](https://github.blog/2021-08-24-github-cli-2-0-includes-extensions/)

# Install

```sh
gh extension install otiai10/gh-dependents
```

# Usage

```sh
gh dependents otiai10/lookpath
# gh dependents {user}/{repo}
```

# Advanced Usage

```sh
gh dependents \
    -v \                     # Show verbose log
    -page=2 \                # Only crawl 2 pages
    -after=MjM1ODQzNDY1NzY \ # Only crawl after specific hash
    -sort=fork \             # Sort output by num of forks
    -pretty \                # Output in pretty format
    otiai10/lookpath

# For more information, hit
gh dependents -h
```

# How it works

- This command just crawls `/network/dependents` page of your repository.

# Usage as a Go library

```go
package main

import (
    "github.com/otiai10/gh-dependents/ghdeps"
)

func main() {
    crawler := ghdeps.NewCrawler("otiai10/lookpath")
	err := crawler.All()
    if err != nil {
        log.Fatalln(err)
    }
    fmt.Printf("%+v\n", crawler.Dependents)
}
```

For more information, check https://pkg.go.dev/github.com/otiai10/gh-dependents/ghdeps

# Issues and Feature Request

- https://github.com/otiai10/gh-dependents/issues
