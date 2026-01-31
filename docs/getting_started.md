<img src="assets/images/banner.png" width="300px">

#

## Getting Started

Instructions to run the project locally for development and testing.

See [deployment notes](deployment.md) for deploying to production.

### Prerequisites

#### Source

    $ git clone git@github.com:frndchagas/standup-raven.git

#### Go 1.25+

    https://golang.org/doc/install

#### Bun

    https://bun.sh

#### Make

On macOS, install XCode command line tools:

    $ xcode-select --install

On Ubuntu:

    $ sudo apt-get install build-essential

#### mmctl

Only needed for `make deploy`. Already included with the Mattermost server or can be installed separately:

    https://docs.mattermost.com/manage/mmctl-command-line-tool.html

### Mattermost Server (dev)

The plugin requires a running Mattermost instance (v9.0.0+). The server repository is at `../mattermost/server/`.

#### Start infrastructure + server

```bash
cd ../mattermost/server
make run-server          # starts Docker (postgres, redis, minio...) + server on port 8065
```

#### Other useful commands

```bash
make run                 # server + webapp
make stop-server         # stop the server
make restart-server      # restart with hot reload
make test-data           # start with sample data (sysadmin / Sys@dmin-sample1)
```

#### Access

    http://localhost:8065

#### Docker Services

| Service    | Port |
|------------|------|
| PostgreSQL | 5433 |
| Redis      | 6379 |
| MinIO      | 9000 |
| Grafana    | 3001 |
| Prometheus | 9090 |
| Inbucket   | 9001 |

### Building

```bash
make dist
```

Generates a `.tar.gz` package in `dist/` containing binaries for all 5 platforms:

| Platform      | Binary |
|---------------|--------|
| Linux amd64   | `server/plugin-linux-amd64` |
| Linux arm64   | `server/plugin-linux-arm64` |
| macOS amd64   | `server/plugin-darwin-amd64` |
| macOS arm64   | `server/plugin-darwin-arm64` |
| Windows amd64 | `server/plugin-windows-amd64.exe` |

### Local Deploy

```bash
make deploy    # installs via mmctl on local Mattermost
```

### Running Tests

```bash
make test
```

### Style Check

```bash
make check-style           # server + webapp
make check-style-server    # Go only (golangci-lint)
make check-style-webapp    # JS only (eslint + stylelint)
```

### Tech Stack

| Component | Technology |
|-----------|------------|
| Server    | Go 1.25, Mattermost Plugin API (server/public v0.1.22) |
| Webapp    | React 18, Bun (bundler + package manager) |
| CI/CD     | GitHub Actions |
| Deploy    | mmctl |
