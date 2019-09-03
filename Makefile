
# Simple Makefile for ercole agent

DESTDIR=build

all: ercole-agent-virtualization

default: ercole-agent-virtualization

clean:
	rm -rf ercole-agent-virtualization build

ercole-agent-virtualization:
	go build -o ercole-agent-virtualization

install: all install-fetchers install-bin install-bin install-config

install-fetchers:
	install -d $(DESTDIR)/fetch
	cp -rp fetch/* $(DESTDIR)/fetch

install-bin:
	install -m 755 ercole-agent-virtualization $(DESTDIR)/ercole-agent-virtualization

install-config:
	install -m 644 config.json $(DESTDIR)/config.json
