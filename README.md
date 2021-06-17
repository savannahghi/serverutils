[![Build Status](https://travis-ci.com/savannahghi/server_utils.svg?branch=main)](https://travis-ci.com/savannahghi/server_utils)
[![Maintained](https://img.shields.io/badge/Maintained-Actively-informational.svg?style=for-the-badge)](https://shields.io/)
# Server Utils Library
serverutils are utilities used by several backend services. These includes starting up and running the server.

### Installing it
serverutils is compatible with modern Go releases in module mode, with Go installed:

```bash
go get -u github.com/savannahghi/serverutils

```
will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/savannahghi/serverutils"

```
and run `go get` without parameters.

The package name is `serverutils`


### Developing

The default branch library is `main`

We try to follow semantic versioning ( <https://semver.org/> ). For that reason,
every major, minor and point release should be _tagged_.

```
git tag -m "v0.0.1" "v0.0.1"
git push --tags
```

Continuous integration tests *must* pass on Travis CI. Our coverage threshold
is 90% i.e you *must* keep coverage above 90%.


## Environment variables

In order to run tests, you need to have an `env.sh` file similar to this one:

```bash
# Application settings
export DEBUG=true
export IS_RUNNING_TESTS=true
export SENTRY_DSN=<a Sentry Data Source Name>
export GOOGLE_CLOUD_PROJECT="Google Cloud project id"
```

This file *must not* be committed to version control.

It is important to _export_ the environment variables. If they are not exported,
they will not be visible to child processes e.g `go test ./...`.

These environment variables should also be set up on Travis CI environment variable section.

## Contributing ##
I would like to cover the entire GitHub API and contributions are of course always welcome. The
calling pattern is pretty well established, so adding new methods is relatively
straightforward. See [`CONTRIBUTING.md`](CONTRIBUTING.md) for details.

## Versioning ##

In general, serverutils follows [semver](https://semver.org/) as closely as we
can for tagging releases of the package. For self-contained libraries, the
application of semantic versioning is relatively straightforward and generally
understood. We've adopted the following
versioning policy:

* We increment the **major version** with any incompatible change to
	non-preview functionality, including changes to the exported Go API surface
	or behavior of the API.
* We increment the **minor version** with any backwards-compatible changes to
	functionality, as well as any changes to preview functionality in the GitHub
	API. GitHub makes no guarantee about the stability of preview functionality,
	so neither do we consider it a stable part of the go-github API.
* We increment the **patch version** with any backwards-compatible bug fixes.

## License ##

This library is distributed under the MIT license found in the [LICENSE](./LICENSE)
file.