# SIL API Client

[![pipeline status](https://gitlab.slade360emr.com/go/api-client/badges/master/pipeline.svg)](https://gitlab.slade360emr.com/go/api-client/-/commits/master)

[![coverage report](https://gitlab.slade360emr.com/go/api-client/badges/master/coverage.svg)](https://gitlab.slade360emr.com/go/api-client/-/commits/master)

This package provides an API client that works with Slade 360 REST APIs and
Slade 360 auth server.

## Installing

### Setting up private Go Modules

Configure GIT to rewrite requests to our Gitlab to occur over SSH:

- `git config --global url."git@gitlab.slade360emr.com:".insteadOf "https://gitlab.slade360emr.com/"`

Add this module to the `GOPRIVATE` list e.g 

```
export GOPRIVATE="gitlab.slade360emr.com/go/api-client"
```

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

## Environment variables

In order to run tests, you need to have an `env.sh` file similar to this one:

```bash
# Application settings
export DEBUG=true
export IS_RUNNING_TESTS=true
export IS_CI=false

# Test API settings
export HOST=<a host>
export API_SCHEME=https
export TOKEN_URL=<a Slade auth server token URL>
export CLIENT_ID=<an OAUTh client ID>
export CLIENT_SECRET=<an OAuth2 client secret>
export USERNAME=<an email to log in with>
export PASSWORD=<a valid password>
export GRANT_TYPE=password
export DEFAULT_WORKSTATION_ID=<an example of a custom header>
```

This file *must not* be committed to version control.

It is important to _export_ the environment variables. If they are not exported,
they will not be visible to child processes e.g `go test ./...`.

These environment variables should also be set up on CI e.g at 
https://gitlab.slade360emr.com/go/api-client/-/settings/ci_cd .
