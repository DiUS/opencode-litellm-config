BINARY := opencode-litellm-config
DIST := dist

PLATFORMS := linux/amd64 darwin/amd64 darwin/arm64 windows/amd64

# Get latest tag, default to 0.0.0 if none
LATEST_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0")
MAJOR := $(word 1,$(subst ., ,$(LATEST_TAG)))
MINOR := $(word 2,$(subst ., ,$(LATEST_TAG)))
PATCH := $(word 3,$(subst ., ,$(LATEST_TAG)))

VERSION ?= $(LATEST_TAG)

.PHONY: build-all checksums release release-major release-minor release-patch clean

build-all: clean
	@mkdir -p $(DIST)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output=$(DIST)/$(BINARY)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then output=$$output.exe; fi; \
		echo "Building $$os/$$arch..."; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -ldflags "-X main.version=$(VERSION)" -o $$output .; \
	done

checksums: build-all
	@cd $(DIST) && sha256sum * > checksums.txt
	@echo "Checksums written to $(DIST)/checksums.txt"

release:
	@if [ "$(VERSION)" = "$(LATEST_TAG)" ] && git rev-parse v$(VERSION) >/dev/null 2>&1; then \
		echo "Error: tag v$(VERSION) already exists. Use VERSION=x.x.x or release-major/minor/patch"; \
		exit 1; \
	fi
	$(MAKE) checksums VERSION=$(VERSION)
	git tag v$(VERSION)
	git push origin v$(VERSION)
	gh release create v$(VERSION) $(DIST)/* --generate-notes --title "v$(VERSION)"

release-major:
	$(MAKE) release VERSION=$(shell echo $$(($(MAJOR)+1))).0.0

release-minor:
	$(MAKE) release VERSION=$(MAJOR).$(shell echo $$(($(MINOR)+1))).0

release-patch:
	$(MAKE) release VERSION=$(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))

clean:
	rm -rf $(DIST)
