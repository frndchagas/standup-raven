<img src="assets/images/banner.png" width="300px">

#

## Getting Started

Instructions to run the project locally for development and testing.

See [deployment notes](deployment.md) for deploying to production.

### Prerequisites

#### Source

    $ git clone git@github.com:standup-raven/standup-raven.git

#### Go 1.22+

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

The plugin requires a running Mattermost instance. The server repository is at `../mattermost/server/`.

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

Generates `.tar.gz` packages in `dist/` for 5 platforms:

| Platform      | File |
|---------------|------|
| Linux amd64   | `mattermost-plugin-standup-raven-vX.Y.Z-linux-amd64.tar.gz` |
| Linux arm64   | `mattermost-plugin-standup-raven-vX.Y.Z-linux-arm64.tar.gz` |
| macOS amd64   | `mattermost-plugin-standup-raven-vX.Y.Z-darwin-amd64.tar.gz` |
| macOS arm64   | `mattermost-plugin-standup-raven-vX.Y.Z-darwin-arm64.tar.gz` |
| Windows amd64 | `mattermost-plugin-standup-raven-vX.Y.Z-windows-amd64.tar.gz` |

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
| Server    | Go 1.22, Mattermost Plugin API v8+ |
| Webapp    | React 17, Bun (bundler + package manager) |
| CI/CD     | GitHub Actions |
| Deploy    | mmctl |
