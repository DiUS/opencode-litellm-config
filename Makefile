BINARY := opencode-litellm-config
VERSION ?= $(shell git describe --tags --always)
DIST := dist

PLATFORMS := linux/amd64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: build-all checksums release clean

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
ifndef VERSION
	$(error VERSION is required. Usage: make release VERSION=0.1.0)
endif
	@if git rev-parse v$(VERSION) >/dev/null 2>&1; then \
		echo "Error: tag v$(VERSION) already exists"; \
		exit 1; \
	fi
	$(MAKE) checksums VERSION=$(VERSION)
	git tag v$(VERSION)
	git push origin v$(VERSION)
	gh release create v$(VERSION) $(DIST)/* --generate-notes --title "v$(VERSION)"

clean:
	rm -rf $(DIST)
