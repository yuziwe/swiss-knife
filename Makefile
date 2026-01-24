# Install location
prefix ?= /usr/local/bin
# Target
target := terminal-translator

terminal-translator:
	go get -C terminal-translator/
	go build -C terminal-translator/ -o ts main.go

install: $(target)
	cp terminal-translator/ts $(prefix)

help:
	@echo "...... make"
	@echo "...... make terminal-translator"
	@echo "...... make install"
	@echo "...... make prefix=<location> install"
	@echo "...... make help"
	@echo "...... make clean"

clean:
	@rm -f terminal-translator/ts

.PHONY: install help clean terminal-translator

