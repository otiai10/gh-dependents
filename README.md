# gh-dependents

[![Actions Status](https://github.com/otiai10/gh-dependents/workflows/Go/badge.svg)](https://github.com/otiai10/gh-dependents/actions)
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
gh dependents -v -t=json otiai10/lookpath
# -v to show verbose log
# -t=json to output in JSON format template
```

# How it works

- This command just crawls `/network/dependents` page of your repository.

# Issues and Feature Request

- https://github.com/otiai10/gh-dependents/issues
