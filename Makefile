GOFLAGS ?= $(GOFLAGS:)
MANIFEST_FILE ?= plugin.json

define GetPluginId
$(shell node -p "require('./plugin.json').id")
endef

define GetPluginVersion
$(shell node -p "'v' + require('./plugin.json').version")
endef

define AddTimeZoneOptions
$(shell node -e \
"let fs = require('fs'); \
try { \
	let manifest = fs.readFileSync('plugin.json', 'utf8'); \
	manifest = JSON.parse(manifest); \
	let timezones = fs.readFileSync('timezones.json', 'utf8'); \
	timezones = JSON.parse(timezones); \
	manifest.settings_schema.settings[0].options=timezones; \
	let json = JSON.stringify(manifest, null, 2); \
	fs.writeFileSync('plugin.json', json, 'utf8'); \
} catch (err) { \
	console.log(err); \
};")
endef

define RemoveTimeZoneOptions
$(shell node -e \
"let fs = require('fs'); \
try { \
	let manifest = fs.readFileSync('plugin.json', 'utf8'); \
	manifest = JSON.parse(manifest); \
	manifest.settings_schema.settings[0].options=[]; \
	let json = JSON.stringify(manifest, null, 2); \
	fs.writeFileSync('plugin.json', json, 'utf8'); \
} catch (err) { \
	console.log(err); \
};")
endef

define UpdateServerHash
git ls-files ./server | xargs shasum -a 256 | cut -d" " -f1 | shasum -a 256 | cut -d" " -f1 > server.sha
endef

define UpdateWebappHash
git ls-files ./webapp | xargs shasum -a 256 | cut -d" " -f1 | shasum -a 256 | cut -d" " -f1 > webapp.sha
endef

PLUGINNAME=$(call GetPluginId)
PLUGINVERSION=$(call GetPluginVersion)
PACKAGENAME=mattermost-plugin-$(PLUGINNAME)-$(PLUGINVERSION)

LDFLAGS=-ldflags "-X 'main.PluginVersion=$(PLUGINVERSION)' -X 'main.SentryServerDSN=$(SERVER_DSN)' -X 'main.SentryWebappDSN=$(WEBAPP_DSN)' -X 'main.EncodedPluginIcon=data:image/svg+xml;base64,$(shell base64 < webapp/src/assets/images/logo.svg | tr -d '\n')'"
GCFLAGS=-gcflags 'all=-N -l'

# All target platforms
PLATFORMS=linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64

.PHONY: default build test run clean stop check-style dist fix-style release deploy

.SILENT: default build test run clean stop check-style dist fix-style release

default: check-style test dist

## Style checking

check-style: check-style-server check-style-webapp

check-style-webapp: .webinstall
	echo Checking for style guide compliance
	cd webapp && bun run lintjs
	cd webapp && bun run lintstyle

check-style-server:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed. Please see https://github.com/golangci/golangci-lint#install for installation instructions."; \
		exit 1; \
	fi
	echo Running golangci-lint
	golangci-lint run ./server/...

fix-style: fix-style-server fix-style-webapp

fix-style-server:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed. Please see https://github.com/golangci/golangci-lint#install for installation instructions."; \
		exit 1; \
	fi
	echo Running golangci-lint --fix
	golangci-lint run --fix ./server/...

fix-style-webapp:
	cd webapp && bun run fixjs
	cd webapp && bun run fixstyle

## Testing

test-server: vendor
	echo Running server tests
	GOTOOLCHAIN=go1.24.11 go test -gcflags='all=-l' -v -coverprofile=coverage.txt ./...

test: test-server

coverage: test-server
	go tool cover -html=coverage.txt -o coverage.html

## Dependencies

.webinstall: webapp/bun.lock
	echo Getting webapp dependencies
	cd webapp && bun install

vendor: go.sum
	echo "Downloading server dependencies"
	go mod download

## Build targets

build-server-%:
	$(eval GOOS=$(firstword $(subst -, ,$*)))
	$(eval GOARCH=$(lastword $(subst -, ,$*)))
	$(eval EXT=$(if $(filter windows,$(GOOS)),.exe,))
	@echo "Building server for $(GOOS)/$(GOARCH)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) $(GCFLAGS) -o dist/intermediate/plugin_$(GOOS)_$(GOARCH)$(EXT) ./server

build-server: $(addprefix build-server-,$(PLATFORMS))

build-webapp: .webinstall
	cd webapp && bun run build
	mkdir -p dist/$(PLUGINNAME)/webapp
	cp -r webapp/dist/* dist/$(PLUGINNAME)/webapp/

## Package & Distribution

package: build-server build-webapp
	@echo "Packaging plugin"
	mkdir -p dist/$(PLUGINNAME)/server
	$(call AddTimeZoneOptions)
	cp plugin.json dist/$(PLUGINNAME)/
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d- -f1); \
		GOARCH=$$(echo $$platform | cut -d- -f2); \
		EXT=""; \
		if [ "$$GOOS" = "windows" ]; then EXT=".exe"; fi; \
		cp dist/intermediate/plugin_$${GOOS}_$${GOARCH}$${EXT} dist/$(PLUGINNAME)/server/plugin-$$platform$${EXT}; \
	done
	cd dist && tar -zcf $(PACKAGENAME).tar.gz $(PLUGINNAME)/*
	$(call RemoveTimeZoneOptions)
	@echo "Built: dist/$(PACKAGENAME).tar.gz"

dist: vendor .webinstall package
	@echo "Distribution packages ready"

## Deploy

deploy:
	@echo "Installing plugin via mmctl"
	@TARBALL="dist/$(PACKAGENAME).tar.gz"; \
	if [ ! -f "$$TARBALL" ]; then \
		echo "Error: $$TARBALL not found. Run 'make dist' first."; \
		exit 1; \
	fi; \
	mmctl plugin add "$$TARBALL" --local && \
	mmctl plugin enable $(PLUGINNAME) --local && \
	echo "Plugin deployed successfully"

## Clean

clean:
	@echo Cleaning plugin
	rm -rf dist/
	rm -rf webapp/node_modules
	rm -rf webapp/.npminstall
	rm -f server.sha server.old.sha
	rm -f webapp.sha webapp.old.sha

## Lint

lint: check-style
