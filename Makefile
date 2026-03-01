# Install location
prefix ?= /usr/local/bin
# Target
target := terminal-translator

terminal-translator:
	go get -C terminal-translator/
	go build -C terminal-translator/ -o ts main.go

openwrt-ruleset-update:
	chmod +x openwrt-ruleset-update/update.fish 
	./openwrt-ruleset-update/update.fish openwrt-ruleset-update/basic-ruleset

install: $(target)
	cp terminal-translator/ts $(prefix)

help:
	@echo "...... make"
	@echo "...... make terminal-translator"
	@echo "...... make openwrt-ruleset-update"
	@echo "...... make install"
	@echo "...... make prefix=<location> install"
	@echo "...... make help"
	@echo "...... make clean"

clean:
	@rm -f terminal-translator/ts
	@rm -rf geoip geosite geoip_url geosite_url

.PHONY: install help clean terminal-translator openwrt-ruleset-update

