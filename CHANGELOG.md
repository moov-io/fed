## v0.7.0 (Released 2021-05-19)

ADDITIONS

- Read `DOWNLOAD_DIRECTORY` environment variable for storing downloaded files.

IMPROVEMENTS

- search: rank results based on fuzzy score rather than name, offer exect routing number matching

BUG FIXES

- Fix file download errors in Docker images
- De-duplicate search results, improve performance

BUILD

- build(deps): bump rexml from 3.2.4 to 3.2.5 in /docs

## v0.6.0 (Released 2021-04-14)

ADDITIONS

- cmd/server: download files if env vars are populated (`FRB_ROUTING_NUMBER` and `FRB_DOWNLOAD_CODE`)

BUILD

- fix(deps): update module github.com/clearbit/clearbit-go to v1.0.1

## v0.5.3 (Released 2021-02-23)

IMPROVEMENTS

- chore(deps): update golang docker tag to v1.16

## v0.5.2 (Released 2021-01-22)

IMPROVEMENTS

- chore(deps): update github.com/xrash/smetrics commit hash to 89a2a8a

BUG FIXES

- build: fixup for OpenShift image running

BUILD

- chore(deps): update golang docker tag to v1.15
- chore(deps): update module gorilla/mux to v1.8.0

## v0.5.1 (Released 2020-07-07)

BUULD

- build: add OpenShift [`quay.io/moov/fed`](https://quay.io/repository/moov/fed) Docker image
- build: convert to Actions from TravisCI
- chore(deps): update module prometheus/client_golang to v1.7.0
- chore(deps): upgrade github.com/gorilla/websocket to v1.4.2

## v0.5.0 (Released 2020-04-14)

ADDITIONS

- ach: support reading input files in the official JSON format
- wire: read official JSON data files

BUILD

- wire: read official JSON data files

## v0.4.3 (Released 2020-03-16)

BUILD

- Fix `make dist` on Windows

## v0.4.2 (Released 2020-03-16)

ADDITIONS

- build: release windows binary

IMPROVEMENTS

- api: use shared Error model
- docs: clarify included data files are old

BUILD

- chore(deps): update golang docker tag to v1.14
- Update module prometheus/client_golang to v1.3.0
- chore(deps): update golang.org/x/oauth2 commit hash to bf48bf1
- build: run sonatype-nexus-community/nancy in CI

## v0.4.1 (Released 2019-12-17)

IMPROVEMENTS

- build: slim down final image, run as moov user

BUILD

- build: test docker image in CI
- Update module prometheus/client_golang to v1.2.1
- build: upgrade openapi-generator to 4.2.2

## v0.4.0 (Released 2019-10-07)

BUG FIXES

- changing ach participant model so AchLocation isn't a list (#68)
- cmd/server: return after marshaling errNoSearchParams

IMPROVEMENTS

- cmd/fedtest: initial binary to perform ACH and Wire searches
- cmd/server: log x-request-id and x-user-id HTTP headers

BUILD

- update module moov-io/base to v0.10.0
- build: upgrade to Go 1.13 and Debian 10

## v0.3.0 (Released 2019-08-16)

BREAKING CHANGES

We've renamed all OpenAPI fields like `Id` to `ID` to be consistent with Go's style.

ADDITIONS

- add environment variables to override command line flags (`LOG_FORMAT`, `HTTP_BIND_ADDRESS`, `HTTP_ADMIN_BIND_ADDRESS`)
- cmd/server: bind HTTP server with TLS if HTTPS_* variables are defined

IMPROVEMENTS

- docs: update docs.moov.io links after design refresh
- docs: link to app specific docs.moov.io page
- cmd/server: quit with an exit code of 1 on missing data files

BUILD

- chore(deps): update module prometheus/client_golang to v1.1.0
- build: download tools used by TravisCI instead of installing them

## v0.2.0 (Released 2019-06-19)

BUILD

- Ship old example FED data files in the Docker image. Production deployments need to replace these with updated files from their Financial Institution.

## v0.1.x (Released 2019-03-06)

BUG FIXES

- Fix automated build steps and Docker setup

ADDITIONS

- Added environmental variables for data filepaths

## v0.1.0 (Released 2019-03-06)

- Initial release
