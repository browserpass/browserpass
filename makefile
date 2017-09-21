SHELL := /usr/bin/env bash
CHROME := $(shell which google-chrome 2>/dev/null || which google-chrome-stable 2>/dev/null || which chromium 2>/dev/null || which chromium-browser 2>/dev/null || which chrome 2>/dev/null)
PEM := $(shell find . -name "chrome-browserpass.pem")
JS_OUTPUT := chrome/script.js chrome/inject.js

all: deps js browserpass

.PHONY: crx
crx:
ifeq ($(PEM),./chrome-browserpass.pem)
	"$(CHROME)" --disable-gpu --pack-extension=./chrome --pack-extension-key=$(PEM)
else
	"$(CHROME)" --disable-gpu --pack-extension=./chrome
	rm ./chrome.pem
endif
	mv chrome.crx chrome-browserpass.crx

.PHONY: js
js: $(JS_OUTPUT)
	cp chrome/host.json chrome-host.json
	cp firefox/host.json firefox-host.json
	cp chrome/{*.html,*.css,*.js,*.png,*.svg} firefox/

chrome/script.js: chrome/script.browserify.js
	browserify chrome/script.browserify.js -o chrome/script.js

chrome/inject.js: chrome/inject.browserify.js
	browserify chrome/inject.browserify.js -o chrome/inject.js

browserpass: cmd/browserpass/ pass/ browserpass.go
	go build -o $@ ./cmd/browserpass

browserpass-linux64: cmd/browserpass/ pass/ browserpass.go
	env GOOS=linux GOARCH=amd64 go build -o $@ ./cmd/browserpass

browserpass-windows64: cmd/browserpass/ pass/ browserpass.go
	env GOOS=windows GOARCH=amd64 go build -o $@.exe ./cmd/browserpass

browserpass-darwinx64: cmd/browserpass/ pass/ browserpass.go
	env GOOS=darwin GOARCH=amd64 go build -o $@ ./cmd/browserpass

browserpass-openbsd64: cmd/browserpass/ pass/ browserpass.go
	env GOOS=openbsd GOARCH=amd64 go build -o $@ ./cmd/browserpass

browserpass-freebsd64: cmd/browserpass/ pass/ browserpass.go
	env GOOS=freebsd GOARCH=amd64 go build -o $@ ./cmd/browserpass

clean:
	rm -f browserpass
	rm -f browserpass-*
	rm -rf release
	git clean -fdx chrome/
	git clean -fdx firefox/
	rm -f *.crx
	rm -f *-host.json

sign: release
	for file in release/*; do \
		gpg --detach-sign "$$file"; \
	done

deps:
	yarn
	dep ensure

tarball: clean deps js
	rm -rf /tmp/browserpass /tmp/browserpass-src.tar.gz
	cp -r ../browserpass /tmp/browserpass
	rm -rf /tmp/browserpass/.git
	find /tmp/browserpass -name "*.pem" -type f -delete
	(cd /tmp && tar -czf /tmp/browserpass-src.tar.gz browserpass)
	mkdir -p release
	cp /tmp/browserpass-src.tar.gz release/

.PHONY: release js crx
release: clean deps js tarball crx browserpass-linux64 browserpass-darwinx64 browserpass-openbsd64 browserpass-freebsd64 browserpass-windows64
	mkdir -p release
	cp chrome-browserpass.crx release/
	zip -jFS "release/chrome" chrome/* chrome-browserpass.crx
	zip -jFS "release/firefox" firefox/*
	zip -FS "release/browserpass-linux64" browserpass-linux64 *-host.json chrome-browserpass.crx install.sh README.md LICENSE
	zip -FS "release/browserpass-darwinx64" browserpass-darwinx64 *-host.json chrome-browserpass.crx install.sh README.md LICENSE
	zip -FS "release/browserpass-openbsd64" browserpass-openbsd64 *-host.json chrome-browserpass.crx install.sh README.md LICENSE
	zip -FS "release/browserpass-freebsd64" browserpass-freebsd64 *-host.json chrome-browserpass.crx install.sh README.md LICENSE
	zip -FS "release/browserpass-windows64" browserpass-windows64.exe *-host.json chrome-browserpass.crx *.ps1 README.md LICENSE
