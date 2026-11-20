---
layout: page
title: Configuration
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Configuration

| Environmental variable | Description | Default |
|-----|-----|-----|
| `FEDACH_DATA_PATH` | Filepath to FedACH data file | `./data/FedACHdir.txt` |
| `FEDWIRE_DATA_PATH` | Filepath to Fedwire data file | `./data/fpddir.txt` |
| `INITIAL_DATA_DIRECTORY` | Directory of files to use instead of downloading or the `*_DATA_PATH` variables. Fed looks for `FedACHdir.txt`, `fedachdir.json`, `fedach.txt`, or `fedach.json` (ACH) and `fpddir.txt`, `fpddir.json`, `fedwire.txt`, or `fedwire.json` (Wire). | Empty |
| `FRB_ROUTING_NUMBER` | Routing number (org id) used to download files from Federal Reserve eServices. Must be set with `FRB_DOWNLOAD_CODE`. | Empty |
| `FRB_DOWNLOAD_CODE` | Download code issued for E-Payments Routing Directory access. Must be set with `FRB_ROUTING_NUMBER`. A missing or incorrect code typically returns HTTP 200 with `"response":{"code":202}` and an empty participant list; the server logs the FRB response code and exits. | Empty |
| `FRB_DOWNLOAD_URL_TEMPLATE` | URL template for the FRB download. Fed replaces `%s` with `fedach` or `fedwire`. | `https://frbservices.org/EPaymentsDirectory/directories/%s?format=json` |
| `LOG_FORMAT` | Format for logging lines to be written as. | Options: `json`, `plain` — Default: `plain` |
| `HTTP_BIND_ADDRESS` | Address for Fed to bind its HTTP server on. This overrides the command-line flag `-http.addr`. | Default: `:8086` |
| `HTTP_ADMIN_BIND_ADDRESS` | Address for Fed to bind its admin HTTP server on. This overrides the command-line flag `-admin.addr`. | Default: `:9096` |
| `HTTP_ADMIN_PPROF_ENABLED` | Set `true`/`1`/`yes` to enable or `false`/`0`/`no` to disable Go pprof on the admin server. Unset keeps historical pprof on. | Unset (enabled) |
| `HTTP_ADMIN_PPROF_SECRET` | Optional bearer token required for pprof when it is enabled. | Empty |
| `HTTPS_CERT_FILE` | Filepath containing a certificate (or intermediate chain) to be served by the HTTP server. Requires all traffic be over secure HTTP. | Empty |
| `HTTPS_KEY_FILE` | Filepath of a private key matching the leaf certificate from `HTTPS_CERT_FILE`. | Empty |
| `BASE_PATH` | HTTP path to serve the API and web UI from. | `/` |
