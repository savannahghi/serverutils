# SIL Base Library

[![pipeline status](https://gitlab.slade360emr.com/go/base/badges/master/pipeline.svg)](https://gitlab.slade360emr.com/go/base/-/commits/master)

[![coverage report](https://gitlab.slade360emr.com/go/base/badges/master/coverage.svg)](https://gitlab.slade360emr.com/go/base/-/commits/master)

This package provides common utilities for our Go servers e.g an API client
that works with Slade 360 REST APIs and Slade 360 auth server.

## Installing

### Setting up private Go Modules

Configure GIT to rewrite requests to our Gitlab to occur over SSH:

- `git config --global url."git@gitlab.slade360emr.com:".insteadOf "https://gitlab.slade360emr.com/"`

Add this module to the `GOPRIVATE` list e.g

```
export GOPRIVATE="gitlab.slade360emr.com/go/base"
```

If you have SSH for Gitlab configured correctly, this should work. If you run
into problems, see <https://stackoverflow.com/a/45936697> for some more ideas.

An alternative:

```
git config \
  --global \
  url."https://oauth2:${GITLAB_REPOSITORY_PERSONAL_ACCESS_TOKEN}@gitlab.slade360emr.com".insteadOf \
  "https://gitlab.slade360emr.com.com"
```

You can create your personal access token at: <https://gitlab.slade360emr.com/profile/personal_access_tokens> .
This personal access token should be granted the `read_repository` permission.
On CI, this should be set up as a _masked_ environment variable, under the name
`GITLAB_REPOSITORY_PERSONAL_ACCESS_TOKEN`.

### Installing it

To install:

```
go get -u gitlab.slade360emr.com/go/base
```

The package name is `client`.

## Developing

The default branch for these small libraries is `master`.

We try to follow semantic versioning ( <https://semver.org/> ). For that reason,
every major, minor and point release should be _tagged_.

```
git tag -m "v0.0.1" "v0.0.1"
git push --tags
```

Continuous integration tests *must* pass on Gitlab CI. Our coverage threshold
for small libraries is 90% i.e you *must* keep coverage above 90%.

## Environment variables

In order to run tests, you need to have an `env.sh` file similar to this one:

```bash
# Application settings
export DEBUG=true
export IS_RUNNING_TESTS=true
export SENTRY_DSN=<a Sentry Data Source Name>

# Google Cloud credentials
export GOOGLE_APPLICATION_CREDENTIALS="<path to a service account JSON file"
export GOOGLE_CLOUD_PROJECT=bewell-app-ci
export FIREBASE_WEB_API_KEY="<a web API key that corresponds to the project named above>"

# Link shortening
export FIREBASE_DYNAMIC_LINKS_DOMAIN=https://bwlci.page.link
export SERVER_PUBLIC_DOMAIN=https://api-gateway-test.healthcloud.co.ke

# Test API settings
export HOST=erp-api-staging.healthcloud.co.ke
export API_SCHEME=https
export TOKEN_URL=https://auth.healthcloud.co.ke/oauth2/token/
export CLIENT_ID="<a valid OAUTH2 client ID>"
export CLIENT_SECRET="<a valid OAUTH2 client secret>"
export USERNAME="<a valid username on Slade 360 auth server>"
export PASSWORD="<a valid password on Slade 360 auth server>"
export GRANT_TYPE=password
export DEFAULT_WORKSTATION_ID="<a workstation ID from Slade 360 ERP that has been linked to the user above>"
export ROOT_COLLECTION_SUFFIX="testing"

export SERVICE_ENVIRONMENT_SUFFIX="testing"
export JWT_KEY="<a secret key>"
```

This file *must not* be committed to version control.

It is important to _export_ the environment variables. If they are not exported,
they will not be visible to child processes e.g `go test ./...`.

These environment variables should also be set up on CI e.g at
<https://gitlab.slade360emr.com/go/base/-/settings/ci_cd> .
