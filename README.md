# SIL API Client

[![pipeline status](https://gitlab.slade360emr.com/go/api-client/badges/master/pipeline.svg)](https://gitlab.slade360emr.com/go/api-client/-/commits/master)

[![coverage report](https://gitlab.slade360emr.com/go/api-client/badges/master/coverage.svg)](https://gitlab.slade360emr.com/go/api-client/-/commits/master)

This package provides an API client that works with Slade 360 REST APIs and
Slade 360 auth server.

## Installing

### Setting up private Go Modules

Configure GIT to rewrite requests to our Gitlab to occur over SSH:

- `git config --global url."git@gitlab.slade360emr.com:".insteadOf "https://gitlab.slade360emr.com/"`

If you have SSH for Gitlab configured correctly, this should work. If you run
into problems, see https://stackoverflow.com/a/45936697 for some more ideas.

### Installing it

To install:

```
go get -u gitlab.slade360emr.com/go/api-client
```

The package name is `client`.

## Developing

The default branch for these small libraries is `master`. 

We try to follow semantic versioning ( https://semver.org/ ). For that reason,
every major, minor and point release should be _tagged_.

```
git tag -m "v0.0.1" "v0.0.1"
git push --tags
```

Continous integration tests *must* pass on Gitlab CI. Our coverage threshold
for small libraries is 90% i.e you *must* keep coverage above 90%.

# TODO Small CI container
