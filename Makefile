HOSTNAME  = registry.terraform.io
NAMESPACE = JRaver
NAME      = ldap
BINARY    = terraform-provider-$(NAME)
VERSION   = $(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.1.0")

OS   := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
OS_ARCH = $(OS)_$(ARCH)

PLUGIN_DIR = ~/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS_ARCH)

.PHONY: default build install generate-docs

default: install

build:
	go build -o $(BINARY)

install: build
	mkdir -p $(PLUGIN_DIR)
	mv $(BINARY) $(PLUGIN_DIR)
	@echo "Installed $(BINARY) v$(VERSION) for $(OS_ARCH)"

generate-docs:
	tfplugindocs
