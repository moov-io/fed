---
layout: page
title: API configuration
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Configuration settings

| Environmental Variable | Description | Default |
|-----|-----|-----|
| `FEDACH_DATA_PATH` | Filepath to FedACH data file | `./data/FedACHdir.txt` |
| `FEDWIRE_DATA_PATH` | Filepath to Fedwire data file | `./data/fpddir.txt` |
| `LOG_FORMAT` | Format for logging lines to be written as. | Options: `json`, `plain` - Default: `plain` |
| `HTTP_BIND_ADDRESS` | Address for Fed to bind its HTTP server on. This overrides the command-line flag `-http.addr`. | Default: `:8086` |
| `HTTP_ADMIN_BIND_ADDRESS` | Address for Fed to bind its admin HTTP server on. This overrides the command-line flag `-admin.addr`. | Default: `:9096` |
| `HTTPS_CERT_FILE` | Filepath containing a certificate (or intermediate chain) to be served by the HTTP server. Requires all traffic be over secure HTTP. | Empty |
| `HTTPS_KEY_FILE`  | Filepath of a private key matching the leaf certificate from `HTTPS_CERT_FILE`. | Empty |

## Data persistence
By design, Fed  **does not persist** (save) any data about the search queries created. The only storage occurs in memory of the process and upon restart Fed will have no files or data saved. Also, no in-memory encryption of the data is performed.