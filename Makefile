VERSION := $$(make -s show-version)
MAIN = ./cmd/gcssurfer

.PHONY: show-version
show-version:
	@gobump show -r $(MAIN)

.PHONY: test
test:
	go test -v ./...

.PHONY: security
security:
	gosec ./...

.PHONY: release
release:
	gobump up -w "$(MAIN)"
	goxc
	ghr "v$(VERSION)" ./dist/snapshot/

	