SHELL := /usr/bin/env bash

.PHONY: empty
empty:

.PHONY: chrome
chrome:
	google-chrome --pack-extension=./chrome --pack-extension-key=./chrome-browserpass.pem
	mv chrome.crx chrome-browserpass.crx

.PHONY: firefox
firefox:
	cp chrome/{*.html,*.css,*.js,*.png,*.svg} firefox/

.PHONY: js
js: chrome/script.browserify.js
	browserify chrome/script.browserify.js -o chrome/script.js

.PHONY: static-files
static-files: chrome/host.json firefox/host.json
	cp chrome/host.json chrome-host.json
	cp firefox/host.json firefox-host.json

browserpass-linux64: browserpass.go
	env GOOS=linux GOARCH=amd64 go build -o $@

browserpass-darwinx64: browserpass.go
	env GOOS=darwin GOARCH=amd64 go build -o $@

.PHONY: static-files chrome firefox
release: static-files chrome firefox browserpass-linux64 browserpass-darwinx64
	mkdir -p release
	zip -jFS "release/chrome" chrome/* chrome-browserpass.crx
	zip -jFS "release/firefox" firefox/*
	zip -FS "release/browserpass-linux64" browserpass-linux64 *-host.json chrome-browserpass.crx install.sh README.md LICENSE
	zip -FS "release/browserpass-darwinx64" browserpass-darwinx64 *-host.json chrome-browserpass.crx install.sh README.md LICENSE
