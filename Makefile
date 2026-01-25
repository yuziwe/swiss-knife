GOENV=CGO_ENABLED=0
GO=$(GOENV) $(shell which go)
GOFLAGS=-ldflags="-extldflags -static -s -w" -trimpath
PREFIX?=/usr/local/bin
TARGET=terminal-translator

terminal-translator:
	@$(GO) build -C terminal-translator/ $(GOFLAGS) -o ts main.go

install: $(TARGET)
	@cp terminal-translator/ts $(PREFIX) || echo "...... Install failed!"

help:
	@echo "...... make"
	@echo "...... make terminal-translator"
	@echo "...... make install"
	@echo "...... make PREFIX=<location> install"
	@echo "...... make help"
	@echo "...... make clean"

clean:
	@rm -f terminal-translator/ts

.PHONY: install help clean terminal-translator

