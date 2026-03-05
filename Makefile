PREFIX ?= $(HOME)/.local

build:
	go build -o bin/claude-watch .

install: build
	install -d $(PREFIX)/bin
	install bin/claude-watch $(PREFIX)/bin/claude-watch

clean:
	rm -rf bin

.PHONY: build install clean
