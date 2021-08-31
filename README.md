# gh-dependents

[![Actions Status](https://github.com/otiai10/gh-dependents/workflows/Go/badge.svg)](https://github.com/otiai10/gh-dependents/actions)
[![codecov](https://codecov.io/gh/otiai10/gh-dependents/branch/main/graph/badge.svg)](https://codecov.io/gh/otiai10/gh-dependents)

`gh` command extension to see dependents of your repository.

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

# How it works

- This command just crawls `/network/dependents` page of your repository.
