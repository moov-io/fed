## v0.13.0 (Released 2025-05-15)

ADDITIONS

- feat: read files from INITIAL_DATA_DIRECTORY

IMPROVEMENTS

- cmd/server: fail on startup when zero participants are found

BUILD

- chore(deps): update dependency go to v1.24.3 (#338)
- fix(deps): update module github.com/moov-io/base to v0.55.0 (#339)
- fix(deps): update module golang.org/x/oauth2 to v0.30.0 (#336)
- fix(deps): update module golang.org/x/text to v0.25.0 (#337)

## v0.12.1 (Released 2025-04-22)

IMPROVEMENTS

- cmd/server: improve startup logging
- cmd/server: return list stats in search response

BUILD

- build: update dependencies
- chore(deps): update dependency go to v1.24.1 (#319)

## v0.11.1 (Released 2024-03-05)

IMPROVEMENTS

- fix: wire file parse bug, improve logic for picking json vs plaintext parser

## v0.11.0 (Released 2024-02-27)

IMPROVEMENTS

- feat: Enable fedach and fedwire download from proxy
- fix: close xml encoder and remove unneeded panics

BUILD

- build: use latest stable Go release
- fix(deps): update module golang.org/x/oauth2 to v0.16.0
- fix(deps): update module github.com/moov-io/base to v0.48.5

## v0.10.2 (Released 2023-06-15)

IMPROVEMENTS

- Use Go 1.20.x in build, update deps

## v0.10.1 (Released 2023-04-19)

REVERTS

- client: revert openapi-generator back to 4.3.1

## v0.10.0 (Released 2023-04-12)

IMPROVEMENTS

- chore: update openapi-generator from 4.2.2 to 6.5.0
- fix: lowercase ach/wire participants in OpenAPI

BUILD

- build(deps): bump nokogiri from 1.13.10 to 1.14.3 in /docs
- build(deps): bump github.com/moov-io/base from v0.39.0 to v0.40.1

## v0.9.2 (Released 2023-04-07)

IMPROVEMENTS

- fix: check typecast of Logo

BUILD

- build: upgrade golang to 1.20
- fix(deps): update module github.com/moov-io/base to v0.39.0
- bump golang.org/x/net from 0.6.0 to 0.7.0
- build: update github.com/stretchr/testify to v1.8.2

## v0.9.1 (Released 2022-09-29)

BUILD

- build: remove deprecated ioutil functions
- fix(deps): update golang.org/x/oauth2 digest to f213421
- fix(deps): update module github.com/moov-io/base to v0.35.0

## v0.9.0 (Released 2022-08-03)

IMPROVEMENTS

- Remove `DOWNLOAD_DIRECTORY` and store downloaded files in memory.

## v0.8.1 (Released 2022-08-02)

IMPROVEMENTS

- fix: remove achParticipants or wireParticipants from json responses

BUILD

- build: require Go 1.18 and set ReadHeaderTimeout
- fix(deps): update module github.com/moov-io/base to v0.33.0
- fix(deps): update golang.org/x/oauth2 digest to 128564f

## v0.8.0 (Released 2022-05-25)

ADDITIONS

- feat: add clearbit logos in responses when configured
- feat: normalize FRB names prior to clearbit search

IMPROVEMENTS

- fix: improve name search by using cleaned name
- refactor: cleanup duplicate code in search logic

BUILD

- build: update codeql action
- build(deps): bump nokogiri from 1.13.4 to 1.13.6 in /docs

## v0.7.4 (Released 2022-05-18)

BUILD

- build: update base images
- build(deps): bump nokogiri from 1.13.3 to 1.13.4 in /docs
- fix(deps): update golang.org/x/oauth2 digest to 9780585

## v0.7.3 (Released 2022-04-04)

IMPROVEMENTS

- fix: replace deprecated strings.Title

## v0.7.2 (Released 2022-04-04)

BUILD

- build(deps): bump nokogiri from 1.12.5 to 1.13.3 in /docs
- fix(deps): update golang.org/x/oauth2 commit hash to ee48083
- fix(deps): update module github.com/go-kit/kit to v0.12.0
- fix(deps): update module github.com/moov-io/base to v0.28.1
- fix(deps): update module github.com/prometheus/client_golang to v1.12.1

## v0.7.1 (Released 2021-07-16)

BUILD

- build(deps): bump addressable from 2.7.0 to 2.8.0 in /docs
- build(deps): bump nokogiri from 1.11.1 to 1.11.5 in /docs
- fix(deps): update golang.org/x/oauth2 commit hash to d040287
- fix(deps): update module github.com/go-kit/kit to v0.11.0

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
